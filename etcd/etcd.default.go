package etcd

import (
	"context"
	xcontrol "github.com/75912001/xlib/control"
	xerror "github.com/75912001/xlib/error"
	xlog "github.com/75912001/xlib/log"
	xruntime "github.com/75912001/xlib/runtime"
	"github.com/pkg/errors"
	etcdclientv3 "go.etcd.io/etcd/client/v3"
	"runtime/debug"
	"sync"
	"time"
)

type Etcd struct {
	client                        *etcdclientv3.Client
	kv                            etcdclientv3.KV
	lease                         etcdclientv3.Lease
	leaseGrantResponse            *etcdclientv3.LeaseGrantResponse
	leaseKeepAliveResponseChannel <-chan *etcdclientv3.LeaseKeepAliveResponse

	cancelFunc context.CancelFunc
	waitGroup  sync.WaitGroup // Stop 等待信号

	options *Options
}

func NewEtcd(opts ...*Options) *Etcd {
	opt := MergeOptions(opts...)
	err := configure(opt)
	if err != nil {
		xlog.PrintfErr("configure err:%v %v", err, xruntime.Location())
		return nil
	}
	return &Etcd{
		options: opt,
	}
}

// Start 开始
func (p *Etcd) Start(ctx context.Context, value string) error {
	var err error
	p.client, err = etcdclientv3.New(etcdclientv3.Config{
		Endpoints:   p.options.endpoints,
		DialTimeout: *p.options.dialTimeout,
	})
	if err != nil {
		return errors.WithMessagef(err, "etcd new client err. endpoints:%v %v", p.options.endpoints, xruntime.Location())
	}
	// 获得kv api子集
	p.kv = etcdclientv3.NewKV(p.client)
	// 申请一个lease 租约
	p.lease = etcdclientv3.NewLease(p.client)
	// 申请一个ttl秒的租约
	p.leaseGrantResponse, err = p.lease.Grant(context.TODO(), *p.options.ttl)
	if err != nil {
		return errors.WithMessagef(err, "etcd new lease err. endpoints:%v %v", p.options.endpoints, xruntime.Location())
	}
	// 删除
	if p.options.key != nil {
		_, err = p.DelWithPrefix(*p.options.key)
		if err != nil {
			return errors.WithMessage(err, xruntime.Location())
		}
	}
	// 添加
	if 0 < len(value) {
		_, err = p.PutWithLease(*p.options.key, value)
		if err != nil {
			return errors.WithMessage(err, xruntime.Location())
		}
	}
	// etcd-watch
	if err = p.WatchPrefixIntoChan(); err != nil {
		return errors.WithMessage(err, xruntime.Location())
	}
	// etcd-get
	if err = p.GetPrefixIntoChan(); err != nil {
		return errors.WithMessage(err, xruntime.Location())
	}
	return nil
}

// Stop 停止
func (p *Etcd) Stop() error {
	if p.client != nil { // 删除
		if _, err := p.DelWithPrefix(*p.options.key); err != nil {
			xlog.PrintfErr("DelWithPrefix err:%v %v", err, xruntime.Location())
		}
		err := p.client.Close()
		if err != nil {
			return errors.WithMessage(err, xruntime.Location())
		}
		p.client = nil
	}

	if p.cancelFunc != nil {
		p.cancelFunc()
		// 等待 goroutine退出.
		p.waitGroup.Wait()
		p.cancelFunc = nil
	}
	return nil
}

func (p *Etcd) GetKey() string {
	return *p.options.key
}

// 多次重试 Start 和 KeepAlive
func (p *Etcd) retryKeepAlive(ctx context.Context) error {
	const grantLeaseRetryDuration = time.Second * 3 // 授权租约 重试 间隔时长
	xlog.PrintfErr("renewing etcd lease, reconfiguring.grantLeaseMaxRetries:%v, grantLeaseIntervalSecond:%v",
		*p.options.grantLeaseMaxRetries, grantLeaseRetryDuration/time.Second)
	var failedGrantLeaseAttempts = 0
	for {
		if err := p.Start(ctx, ""); err != nil {
			failedGrantLeaseAttempts++
			if *p.options.grantLeaseMaxRetries <= failedGrantLeaseAttempts {
				return errors.WithMessagef(err, "%v exceeded max attempts to renew etcd lease %v %v",
					xruntime.Location(), *p.options.grantLeaseMaxRetries, failedGrantLeaseAttempts)
			}
			xlog.PrintErr("error granting etcd lease, will retry.", err)
			time.Sleep(grantLeaseRetryDuration)
			continue
		} else {
			// 续租
			if err = p.KeepAlive(ctx); err != nil {
				failedGrantLeaseAttempts++
				if *p.options.grantLeaseMaxRetries <= failedGrantLeaseAttempts {
					return errors.WithMessagef(err, "%v exceeded max attempts to renew etcd lease %v %v",
						xruntime.Location(), *p.options.grantLeaseMaxRetries, failedGrantLeaseAttempts)
				}
				xlog.PrintErr("error granting etcd lease, will retry.", err)
				time.Sleep(grantLeaseRetryDuration)
				continue
			} else {
				return nil
			}
		}
	}
}

// KeepAlive 更新租约
func (p *Etcd) KeepAlive(ctx context.Context) error {
	var err error
	p.leaseKeepAliveResponseChannel, err = p.lease.KeepAlive(ctx, p.leaseGrantResponse.ID)
	if err != nil {
		return errors.WithMessage(err, xruntime.Location())
	}
	p.waitGroup.Add(1)
	ctxWithCancel, cancelFunc := context.WithCancel(ctx)
	p.cancelFunc = cancelFunc
	go func(ctx context.Context) {
		defer func() {
			if xruntime.IsRelease() {
				if err := recover(); err != nil {
					xlog.PrintErr(xerror.GoroutinePanic, err, debug.Stack())
				}
			}
			p.waitGroup.Done()
			xlog.PrintInfo(xerror.GoroutineDone)
		}()
		for {
			select {
			case <-ctx.Done():
				xlog.PrintInfo(xerror.GoroutineDone)
				return
			case leaseKeepAliveResponse, ok := <-p.leaseKeepAliveResponseChannel:
				//xlog.PrintInfo(leaseKeepAliveResponse, ok)
				if leaseKeepAliveResponse != nil {
					continue
				}
				if ok {
					continue
				}
				// abnormal
				xlog.PrintErr("etcd lease KeepAlive died, retrying")
				go func(ctx context.Context) {
					defer func() {
						if xruntime.IsRelease() {
							if err := recover(); err != nil {
								xlog.PrintErr(xerror.Retry, xerror.GoroutinePanic, err, debug.Stack())
							}
						}
						xlog.PrintInfo(xerror.Retry, xerror.GoroutineDone)
					}()
					if err := p.Stop(); err != nil {
						xlog.PrintInfo(xerror.Retry, xerror.Fail, err)
						return
					}
					if err := p.retryKeepAlive(ctx); err != nil {
						xlog.PrintErr(xerror.Retry, xerror.Fail, err)
						return
					}
				}(context.TODO())
				return
			}
		}
	}(ctxWithCancel)
	return nil
}

// PutWithLease 将一个键值对放入etcd中 WithLease 带ttl
func (p *Etcd) PutWithLease(key string, value string) (*etcdclientv3.PutResponse, error) {
	putResponse, err := p.kv.Put(context.TODO(), key, value, etcdclientv3.WithLease(p.leaseGrantResponse.ID))
	if err != nil {
		return nil, errors.WithMessagef(err, "etcd put with lease err. key:%v value:%v %v", key, value, xruntime.Location())
	}
	return putResponse, nil
}

// Put 将一个键值对放入etcd中 [不带租约ttl]
//func (p *Etcd) Put(key string, value string) (*etcdclientv3.PutResponse, error) {
//	putResponse, err := p.kv.Put(context.TODO(), key, value)
//	if err != nil {
//		return nil, errors.WithMessage(err, xruntime.Location())
//	}
//	return putResponse, nil
//}

// DelWithPrefix 删除键值 匹配的键值
func (p *Etcd) DelWithPrefix(keyPrefix string) (*etcdclientv3.DeleteResponse, error) {
	deleteResponse, err := p.kv.Delete(context.TODO(), keyPrefix, etcdclientv3.WithPrefix())
	if err != nil {
		return nil, errors.WithMessagef(err, "etcd del with prefix err. keyPrefix:%v %v", keyPrefix, xruntime.Location())
	}
	return deleteResponse, nil
}

//
//// Del 删除键值
//func (p *Mgr) Del(key string) (*clientv3.DeleteResponse, error) {
//	deleteResponse, err := p.kv.Delete(context.TODO(), key)
//	if err != nil {
//		return nil, errors.WithMessage(err, xrutil.GetCodeLocation(1).String())
//	}
//	return deleteResponse, nil
//}
//
//// DelRange 按选项删除范围内的键值
//func (p *Mgr) DelRange(startKeyPrefix string, endKeyPrefix string) (*clientv3.DeleteResponse, error) {
//	opts := []clientv3.OpOption{
//		clientv3.WithPrefix(),
//		clientv3.WithFromKey(),
//		clientv3.WithRange(endKeyPrefix),
//	}
//	deleteResponse, err := p.kv.Delete(context.TODO(), startKeyPrefix, opts...)
//	if err != nil {
//		return nil, errors.WithMessage(err, xrutil.GetCodeLocation(1).String())
//	}
//	return deleteResponse, nil
//}

// WatchPrefix 监视以key为前缀的所有 key value
func (p *Etcd) WatchPrefix(key string) etcdclientv3.WatchChan {
	return p.client.Watch(context.TODO(), key, etcdclientv3.WithPrefix())
}

//
//// Get 检索键
//func (p *Mgr) Get(key string) (*clientv3.GetResponse, error) {
//	getResponse, err := p.kv.Get(context.TODO(), key)
//	if err != nil {
//		return nil, errors.WithMessage(err, xrutil.GetCodeLocation(1).String())
//	}
//	return getResponse, nil
//}

// GetPrefix 查找以key为前缀的所有 key value
func (p *Etcd) GetPrefix(key string) (*etcdclientv3.GetResponse, error) {
	getResponse, err := p.kv.Get(context.TODO(), key, etcdclientv3.WithPrefix())
	if err != nil {
		return nil, errors.WithMessagef(err, "etcd get prefix err. key:%v %v", key, xruntime.Location())
	}
	return getResponse, nil
}

// GetPrefixIntoChan  取得关心的前缀,放入 chan 中
func (p *Etcd) GetPrefixIntoChan() (err error) {
	if p.options.watchKeyPrefix == nil { // 如果没有设置 watchKeyPrefix,则无动作
		return nil
	}
	getResponse, err := p.GetPrefix(*p.options.watchKeyPrefix)
	if err != nil {
		return errors.WithMessage(err, xruntime.Location())
	}
	for _, v := range getResponse.Kvs {
		key := string(v.Key)
		var valueJson *ValueJson
		if len(v.Value) != 0 {
			valueJson = ValueString2Json(string(v.Value))
		}
		_, oldExist := GRegistry.Find(key)
		GRegistry.Update(key, valueJson)
		var cb xcontrol.ICallBack
		if oldExist { // 已存在
			if valueJson == nil { // 删除
				if p.options.DelCallback != nil {
					cb = p.options.DelCallback.Clone(key)
				}
			} else { // 更新
				if p.options.UpdateCallback != nil {
					cb = p.options.UpdateCallback.Clone(key, valueJson)
				}
			}
		} else { // 不存在
			if valueJson == nil { // 删除,但不存在,无动作
				continue
			} else { // 新增
				if p.options.AddCallback != nil {
					cb = p.options.AddCallback.Clone(key, valueJson)
				}
			}
		}
		if cb != nil {
			p.options.iOut.Send(
				&xcontrol.Event{
					ISwitch:   xcontrol.NewSwitchButton(true),
					ICallBack: cb,
				},
			)
		}
	}
	return
}

// WatchPrefixIntoChan 监听key变化,放入 chan 中
func (p *Etcd) WatchPrefixIntoChan() (err error) {
	if p.options.watchKeyPrefix == nil { // 如果没有设置 watchKeyPrefix,则不监听
		return nil
	}
	eventChan := p.WatchPrefix(*p.options.watchKeyPrefix)
	go func() {
		defer func() {
			if xruntime.IsRelease() {
				if err := recover(); err != nil {
					xlog.PrintErr(xerror.GoroutinePanic, err, debug.Stack())
				}
			}
			xlog.PrintInfo(xerror.GoroutineDone)
		}()
		for v := range eventChan {
			kv := v.Events[0].Kv
			key := string(kv.Key)
			Value := string(kv.Value)
			var valueJson *ValueJson
			if len(Value) != 0 {
				valueJson = ValueString2Json(Value)
			}
			_, oldExist := GRegistry.Find(key)
			GRegistry.Update(key, valueJson)
			var cb xcontrol.ICallBack
			if oldExist { // 已存在
				if valueJson == nil { // 删除
					if p.options.DelCallback != nil {
						cb = p.options.DelCallback.Clone(key)
					}
				} else { // 更新
					if p.options.UpdateCallback != nil {
						cb = p.options.UpdateCallback.Clone(key, valueJson)
					}
				}
			} else { // 不存在
				if valueJson == nil { // 删除,但不存在,无动作
					continue
				} else { // 新增
					if p.options.AddCallback != nil {
						cb = p.options.AddCallback.Clone(key, valueJson)
					}
				}
			}
			if cb != nil {
				p.options.iOut.Send(
					&xcontrol.Event{
						ISwitch:   xcontrol.NewSwitchButton(true),
						ICallBack: cb,
					},
				)
			}
		}
	}()
	return
}
