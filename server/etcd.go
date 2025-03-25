package server

import (
	xcontrol "github.com/75912001/xlib/control"
	xetcd "github.com/75912001/xlib/etcd"
)

const ReportIntervalSecondDefault int64 = 30 // etcd-上报时间间隔 秒

// etcdReportFunction etcd-上报
func etcdReportFunction(args ...interface{}) error {
	defaultServer := args[0].(*Server)
	defer func() {
		defaultServer.Timer.AddSecond(
			xcontrol.NewCallBack(etcdReportFunction, defaultServer),
			defaultServer.TimeMgr.ShadowTimestamp()+ReportIntervalSecondDefault,
		)
	}()
	valueJson := &xetcd.ValueJson{
		ServerNet:     defaultServer.ConfigMgr.Config.ServerNet,
		Version:       *defaultServer.ConfigMgr.Config.Base.Version,
		AvailableLoad: defaultServer.AvailableLoad,
		SecondOffset:  0,
	}
	value := xetcd.ValueJson2String(valueJson)
	if _, err := defaultServer.Etcd.Put(defaultServer.EtcdKey, value); err != nil {
		defaultServer.Log.Errorf("etcdReportFunction Put key:%v val:%v err:%v", defaultServer.EtcdKey, value, err)
	}
	return nil
}
