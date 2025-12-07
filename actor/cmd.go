package actor

import (
	"math"
)

const (
	// 用户自定义命令范围: [0,4294959999]
	CustomCommand_Begin uint32 = 0          // 用户自定义命令-起始值
	CustomCommand_End   uint32 = 4294959999 // 用户自定义命令-结束值
	// 系统保留命令范围: [4294960000,4294967295]
	SystemReservedCommand_Begin       uint32 = 4294960000     // 系统保留命令-起始值
	SystemReservedCommand_Stop               = 4294960001     // 停止 {args:无}
	SystemReservedCommand_RemoveChild        = 4294960002     // 移除子 actor {args: [0]:子 actor key}
	SystemReservedCommand_Spawn              = 4294960003     // 创建子 actor {args: [0]:子 actor key, [1]:子 actor 行为函数} {返回: response *Actor[KEY]}
	SystemReservedCommand_GetChild           = 4294960004     // 获取子 actor {args: [0]:子 actor key} {返回: response *Actor[KEY]}
	SystemReservedCommand_End                = math.MaxUint32 // 系统保留命令-结束值
)
