package common

// DisconnectReason 表示断开连接的原因
type DisconnectReason int

const (
	DisconnectReasonUnknown        DisconnectReason = 0   // 未知原因
	DisconnectReasonServerShutdown                  = 100 // 服务器关闭
	DisconnectReasonClientTimeout                   = 200 // 客户端超时
	DisconnectReasonClientLogic                     = 201 // 客户端逻辑
	DisconnectReasonClientShutdown                  = 202 // 客户端关闭
)
