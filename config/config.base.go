package config

import (
	xconfigconstants "github.com/75912001/xlib/config/constants"
	xerror "github.com/75912001/xlib/error"
	xruntime "github.com/75912001/xlib/runtime"
	"github.com/pkg/errors"
	"math"
	"runtime"
)

type Base struct {
	ProjectName                 *string                          `yaml:"projectName"`                 // 项目名称		[default]:constants.ProjectNameDefault
	GroupID                     *uint32                          `yaml:"groupID"`                     // 分组ID
	Name                        *string                          `yaml:"name"`                        // 服务名称
	ServerID                    *uint32                          `yaml:"serverID"`                    // 服务ID
	Version                     *string                          `yaml:"version"`                     // 版本号		[default]: constants.VersionDefault
	PprofHttpPort               *uint16                          `yaml:"pprofHttpPort"`               // pprof性能分析 http端口		[default]: nil 不使用
	GoMaxProcess                *int                             `yaml:"goMaxProcess"`                // [default]: runtime.NumCPU()
	PacketLengthMax             *uint32                          `yaml:"packetLengthMax"`             // bytes,用户 上行 每个包的最大长度		[default]: math.MaxUint32
	SendChannelCapacity         *uint32                          `yaml:"sendChannelCapacity"`         // bytes,每个链接的发送chan容量		[default]: 1000000
	RunMode                     *uint32                          `yaml:"runMode"`                     // 运行模式 [0:release 1:debug]		[default]: 1
	AvailableLoad               *uint32                          `yaml:"availableLoad"`               // 可用资源数		[default]: 1000000
	PacketLimitRecvCntPreSecond *uint32                          `yaml:"packetLimitRecvCntPreSecond"` // 每秒接收包数限制		[default]: math.MaxUint32
	ProcessingMode              *xconfigconstants.ProcessingMode `yaml:"processingMode"`              // 处理-模式 [0:bus 1:actor]		[default]: ProcessingModeBus
}

func (p *Base) ProcessingModeIsActor() bool {
	return *p.ProcessingMode == xconfigconstants.ProcessingModeActor
}

func (p *Base) Configure() error {
	if p.ProjectName == nil {
		p.ProjectName = &xconfigconstants.ProjectNameDefault
	}
	if p.GroupID == nil {
		return errors.WithMessagef(xerror.Config, "groupID is nil. %v", xruntime.Location())
	}
	if p.Name == nil {
		return errors.WithMessagef(xerror.Config, "name is nil. %v", xruntime.Location())
	}
	if p.ServerID == nil {
		return errors.WithMessagef(xerror.Config, "serverID is nil. %v", xruntime.Location())
	}
	if p.Version == nil {
		p.Version = &xconfigconstants.VersionDefault
	}
	if p.GoMaxProcess == nil {
		defaultValue := runtime.NumCPU()
		p.GoMaxProcess = &defaultValue
	}
	if p.PacketLengthMax == nil {
		defaultValue := uint32(math.MaxUint32)
		p.PacketLengthMax = &defaultValue
	}
	if p.SendChannelCapacity == nil {
		defaultValue := uint32(1000000)
		p.SendChannelCapacity = &defaultValue
	}
	if p.RunMode == nil {
		defaultValue := uint32(1)
		p.RunMode = &defaultValue
	}
	if p.AvailableLoad == nil {
		defaultValue := uint32(1000000)
		p.AvailableLoad = &defaultValue
	}
	if p.PacketLimitRecvCntPreSecond == nil {
		defaultValue := uint32(math.MaxUint32)
		p.PacketLimitRecvCntPreSecond = &defaultValue
	}
	if p.ProcessingMode == nil {
		defaultValue := xconfigconstants.ProcessingModeBus
		p.ProcessingMode = &defaultValue
	}
	return nil
}
