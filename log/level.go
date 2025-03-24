package log

// 日志等级
const (
	LevelOff   uint32 = 0 //关闭
	LevelFatal uint32 = 1 //致命
	LevelError uint32 = 2 //错误
	LevelWarn  uint32 = 3 //警告
	LevelInfo  uint32 = 4 //信息
	LevelDebug uint32 = 5 //调试
	LevelTrace uint32 = 6 //跟踪
	LevelOn    uint32 = 7 //全部打开
)

// 等级描述
var levelDesc = []string{
	LevelOff:   "LevelOff", //关闭
	LevelFatal: "FAT",      //致命
	LevelError: "ERR",      //错误
	LevelWarn:  "WAR",      //警告
	LevelInfo:  "INF",      //信息
	LevelDebug: "DEB",      //调试
	LevelTrace: "TRA",      //跟踪
	LevelOn:    "LevelOn",  //全部打开
}
