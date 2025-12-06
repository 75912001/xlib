package server

import (
	xcontrol "github.com/75912001/xlib/control"
	xetcd "github.com/75912001/xlib/etcd"
	xlog "github.com/75912001/xlib/log"
	xtimer "github.com/75912001/xlib/timer"
	"time"
)

const ReportIntervalSecondDefault int64 = 30 // etcd-上报时间间隔 秒

// etcdReportFunction etcd-上报
func etcdReportFunction(args ...any) error {
	defaultServer := args[0].(*Server)
	defer func() {
		xtimer.GTimer.AddSecond(
			xcontrol.NewCallBack(etcdReportFunction, defaultServer),
			time.Now().Unix()+ReportIntervalSecondDefault,
			defaultServer.GetActor(),
		)
	}()
	key := xetcd.GEtcd.GetKey()
	value := defaultServer.genEtcdValue()
	if _, err := xetcd.GEtcd.PutWithLease(key, value); err != nil {
		xlog.GLog.Errorf("etcdReportFunction Put key:%v val:%v err:%v", key, value, err)
	}
	return nil
}
