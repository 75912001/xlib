package constants

const TtlSecondDefault int64 = 33 // etcd 默认 TTL 时间 秒

const WatchMsgTypeServer string = "server"       // etcd watch 消息类型-服务
const WatchMsgTypeCommand string = "command"     // etcd watch 消息类型-命令
const WatchMsgTypeGM string = "gm"               // etcd watch 消息类型-GM
const WatchMsgTypeServerCfg string = "serverCfg" // etcd watch 消息类型-服务配置
