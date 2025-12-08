package subpub

import (
	xcontrol "github.com/75912001/xlib/control"
	xerror "github.com/75912001/xlib/error"
	xmap "github.com/75912001/xlib/map"
	xruntime "github.com/75912001/xlib/runtime"
	"github.com/pkg/errors"
	"sync"
)

// Default 发布订阅器
type Default[KEY ISubPubKey] struct {
	mapMgr *xmap.MapMgr[KEY, []xcontrol.ICallBack]
	mu     sync.RWMutex
}

func NewDefault[KEY ISubPubKey]() *Default[KEY] {
	return &Default[KEY]{
		mapMgr: xmap.NewMapMgr[KEY, []xcontrol.ICallBack](),
	}
}

// Subscribe 订阅
//
//	[❗] 订阅的回调函数, 不可以是闭包函数.
func (p *Default[KEY]) Subscribe(key KEY, onFunction xcontrol.OnFunction) error {
	if onFunction == nil { // 回调为nil
		return errors.WithMessagef(xerror.ParamNotSupport, "onFunction is nil, %v", xruntime.Location())
	}

	p.mu.Lock()
	defer func() {
		p.mu.Unlock()
	}()

	callbacks, ok := p.mapMgr.Find(key)
	if !ok { // 无
		callbacks = []xcontrol.ICallBack{}
	}
	callbacks = append(callbacks, xcontrol.NewCallBack(onFunction))
	p.mapMgr.Add(key, callbacks)
	return nil
}

// Unsubscribe
//
//	[❗] 取消订阅的回调函数, 不可以是闭包函数.
func (p *Default[KEY]) Unsubscribe(key KEY, onFunction xcontrol.OnFunction) error {
	p.mu.Lock()
	defer func() {
		p.mu.Unlock()
	}()

	callbacks, ok := p.mapMgr.Find(key)
	if !ok {
		return nil
	}

	// 创建新的回调列表，排除要取消的回调
	newCallbacks := make([]xcontrol.ICallBack, 0, len(callbacks))
	targetCallback := xcontrol.NewCallBack(onFunction)

	for _, callback := range callbacks {
		if cb, ok := callback.(*xcontrol.CallBack); ok {
			if cb.Equals(targetCallback) { // 相等-取消的订阅
				continue
			}
			newCallbacks = append(newCallbacks, callback)
		} else {
			newCallbacks = append(newCallbacks, callback)
		}
	}

	// 更新或删除回调列表
	if len(newCallbacks) > 0 {
		p.mapMgr.Add(key, newCallbacks)
	} else {
		p.mapMgr.Del(key)
	}
	return nil
}

func (p *Default[KEY]) Publish(key KEY, args ...any) error {
	var returnError error
	p.mu.RLock()
	defer func() {
		p.mu.RUnlock()
	}()

	callbacks, ok := p.mapMgr.Find(key)
	if !ok { // 无
		return nil
	}
	for _, callback := range callbacks {
		if err := callback.Clone(args...).Execute(); err != nil {
			if returnError == nil {
				returnError = errors.WithMessagef(err, "Publish %v %v", key, args)
			} else {
				returnError = errors.WithMessagef(returnError, "%v Publish %v %v error.", err, key, args)
			}
		}
	}
	return returnError
}

func (p *Default[KEY]) Clear() {
	p.mu.Lock()
	defer func() {
		p.mu.Unlock()
	}()

	p.mapMgr.Clear()
}

type ISubPubKey interface {
	int | int32 | uint32 | int64 | uint64 | string
}
