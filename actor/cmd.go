package actor

import (
	"math"
)

type CMD uint32

const (
	// 用户自定义命令范围: [0,4294959999]
	CustomCommand_Begin CMD = 0          // 用户自定义命令-起始值
	CustomCommand_End   CMD = 4294959999 // 用户自定义命令-结束值
	// 系统保留命令范围: [4294960000,4294967295]
	SystemReservedCommand_Begin       CMD = 4294960000     // 系统保留命令-起始值
	SystemReservedCommand_Stop        CMD = 4294960001     // 停止 {args:无}
	SystemReservedCommand_RemoveChild CMD = 4294960002     // 移除子 actor {args: [0]:子 actor key}
	SystemReservedCommand_Spawn       CMD = 4294960003     // 创建子 actor {args: [0]:子 actor key, [1]:子 actor 行为函数} {返回: response *Actor[KEY]}
	SystemReservedCommand_GetChild    CMD = 4294960004     // 获取子 actor {args: [0]:子 actor key} {返回: response *Actor[KEY]}
	SystemReservedCommand_End         CMD = math.MaxUint32 // 系统保留命令-结束值
)

// 是否-用户自定义命令
func isCustomCommand(cmd CMD) bool {
	return cmd <= CustomCommand_End
}

// 是否-系统保留命令
func isSystemReservedCommand(cmd CMD) bool {
	return SystemReservedCommand_Begin <= cmd
}
