package common

import (
	xruntime "github.com/75912001/xlib/runtime"
	"github.com/pkg/errors"
	"net"
	"time"
)

const ServerNetTypeNameTCP = "tcp"
const ServerNetTypeNameKCP = "kcp"
const ServerNetTypeNameWebSocket = "websocket"

// DisconnectReason 表示断开连接的原因
type DisconnectReason int

const (
	DisconnectReasonUnknown DisconnectReason = 0 // 未知原因

	DisconnectReasonClientShutdown DisconnectReason = 1 // 客户端关闭
	DisconnectReasonClientLogic    DisconnectReason = 2 // 客户端逻辑
	DisconnectReasonServerShutdown DisconnectReason = 3 // 服务端关闭
	DisconnectReasonShutdown       DisconnectReason = 4 // 关闭-主动关闭
	DisconnectReasonPeerShutdown   DisconnectReason = 5 // 对端关闭
	// [10000,20000] 留给业务使用
	// ...
)

const (
	EventAddTimeoutDuration = 3 * time.Second // 事件-加入的超时时间
)

// 写超时
//
//	只有超过50%时才更新写截止日期
//	参数:
//		lastTime:上次时间 (可能会更新)
//		thisTime:这次时间
//		writeTimeOutDuration:写超时时长
func UpdateWriteDeadline(conn net.Conn, lastTime *time.Time, thisTime time.Time, writeTimeOutDuration time.Duration) error {
	if (writeTimeOutDuration >> 1) < thisTime.Sub(*lastTime) {
		if err := conn.SetWriteDeadline(thisTime.Add(writeTimeOutDuration)); err != nil {
			return errors.WithMessagef(err, "UpdateWriteDeadline:%v", xruntime.Location())
		}
		*lastTime = thisTime
	}
	return nil
}
