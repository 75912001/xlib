package etcd

import (
	"encoding/json"
	"fmt"
	xlog "github.com/75912001/xlib/log"
	"path"
	"strconv"
)

// e.g.:/${projectName}/${EtcdWatchMsgType}/${groupID}/${serverName}/${serverID}/

// ValueJson etcd 通讯的数据,由服务中的数据生成,定时更新->etcd->服务
type ValueJson struct {
	ServerNet     []*ServerNet `json:"serverNet,omitempty"`     // 有:直接使用. 没有:使用 etcd.ServerNet
	Version       string       `json:"version,omitempty"`       // 有:直接使用. 没有:使用 base.version 生成
	AvailableLoad uint32       `json:"availableLoad,omitempty"` // 剩余可用负载, 可用资源数
	SecondOffset  int64        `json:"secondOffset,omitempty"`  // 服务 时间(秒)偏移量
	GrpcService   *GrpcService `json:"grpcService,omitempty"`   // gRPC 服务信息
}

// Parse 解析key
func Parse(key string) (msgType string, groupID uint32, serviceName string, serviceID uint32) {
	strServiceID := path.Base(key)
	// strServiceID 转换成 serviceID
	if serviceIDU64, err := strconv.ParseUint(strServiceID, 10, 32); err != nil {
		xlog.PrintfErr("Parse err:%v %v %v", key, strServiceID, err)
		return
	} else {
		serviceID = uint32(serviceIDU64)
	}
	key = path.Dir(key)
	key = path.Dir(key)
	serviceName = path.Base(key)
	key = path.Dir(key)
	strGroupID := path.Base(key)
	// strGroupID 转换成 groupID
	if groupIDU64, err := strconv.ParseUint(strGroupID, 10, 32); err != nil {
		xlog.PrintfErr("Parse err:%v %v %v", key, strGroupID, err)
		return
	} else {
		groupID = uint32(groupIDU64)
	}
	key = path.Dir(key)
	msgType = path.Base(key)
	return msgType, groupID, serviceName, serviceID
}

func GenKey(projectName string, etcdWatchMsgType string, groupID uint32, serviceName string, serviceID uint32) string {
	return fmt.Sprintf("/%v/%v/%v/%v/%v/",
		projectName,
		etcdWatchMsgType,
		groupID,
		serviceName,
		serviceID,
	)
}

func GenPrefixKey(projectName string) string {
	return fmt.Sprintf("/%v/", projectName)
}

func ValueJson2String(valueJson *ValueJson) string {
	bytes, err := json.Marshal(valueJson)
	if err != nil {
		xlog.PrintfErr("Error marshaling ValueJson: %v", err)
		return ""
	}
	return string(bytes)
}

func ValueString2Json(value string) *ValueJson {
	var valueJson ValueJson
	err := json.Unmarshal([]byte(value), &valueJson)
	if err != nil {
		xlog.PrintfErr("Error unmarshaling ValueJson: %v %v", value, err)
		return nil
	}
	return &valueJson
}
