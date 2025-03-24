package etcd

import "time"

const TtlSecondDefault int64 = 33            // etcd 默认 TTL 时间 秒
const ReportIntervalSecondDefault int64 = 30 // etcd-上报时间间隔 秒

const WatchMsgTypeServer string = "server"             // etcd watch 消息类型-服务
const WatchMsgTypeCommand string = "command"           // etcd watch 消息类型-命令
const WatchMsgTypeGM string = "gm"                     // etcd watch 消息类型-GM
const WatchMsgTypeServiceBench string = "serviceBench" // etcd watch 消息类型-服务配置

var (
	grantLeaseMaxRetriesDefault = 600             // 授权租约 最大 重试次数
	dialTimeoutDefault          = time.Second * 5 // dialTimeout is the timeout for failing to establish a connection.
)
