// 服务配置
// 服务配置文件, 用于配置服务的基本信息.
// 该配置文件与可执行程序在同一目录下.

package config

import (
	"fmt"
	xcommon "github.com/75912001/xlib/common"
	xerror "github.com/75912001/xlib/error"
	xetcd "github.com/75912001/xlib/etcd"
	xnetcommon "github.com/75912001/xlib/net/common"
	xruntime "github.com/75912001/xlib/runtime"
	xtimer "github.com/75912001/xlib/timer"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
	"path/filepath"
	"runtime"
	"time"
)

// 配置-主项,用户服务的基本配置

type Mgr struct {
	Root   rootYaml
	Config configYaml
}

type rootYaml struct {
	Etcd Etcd `yaml:"etcd"`
}

func (p *rootYaml) Parse(strYaml string) error {
	if err := yaml.Unmarshal([]byte(strYaml), &p); err != nil {
		return errors.WithMessage(err, xruntime.Location())
	}
	if p.Etcd.TTL == nil {
		defaultValue := xetcd.TtlSecondDefault
		p.Etcd.TTL = &defaultValue
	}
	return nil
}

type Etcd struct {
	Addrs []string `yaml:"addrs"` // etcd地址
	TTL   *int64   `yaml:"ttl"`   // ttl 秒 [default]: xetcd.TtlSecondDefault 秒, e.g.:系统每10秒续约一次,该参数至少为11秒
}

type configYaml struct {
	Base      Base                 `yaml:"base"`
	Timer     Timer                `yaml:"timer"`
	ServerNet []*xcommon.ServerNet `yaml:"serverNet"`
}

func (p *configYaml) Parse(yamlString string) error {
	err := yaml.Unmarshal([]byte(yamlString), p)
	if err != nil {
		return errors.WithMessage(err, xruntime.Location())
	}
	if p.Base.ProjectName == nil {
		return errors.WithMessage(err, xruntime.Location())
	}
	if p.Base.Version == nil {
		return errors.WithMessage(err, xruntime.Location())
	}
	if p.Base.LogLevel == nil {
		return errors.WithMessage(err, xruntime.Location())
	}
	if p.Base.LogAbsPath == nil {
		executablePath, err := xruntime.GetExecutablePath()
		if err != nil {
			return errors.WithMessage(err, xruntime.Location())
		}
		executablePath = filepath.Join(executablePath, "log")
		p.Base.LogAbsPath = &executablePath
	}
	if p.Base.GoMaxProcess == nil {
		defaultValue := runtime.NumCPU()
		p.Base.GoMaxProcess = &defaultValue
	}
	if p.Base.BusChannelCapacity == nil {
		return errors.WithMessage(err, xruntime.Location())
	}
	if p.Base.PacketLengthMax == nil {
		return errors.WithMessage(err, xruntime.Location())
	}
	if p.Base.SendChannelCapacity == nil {
		return errors.WithMessage(err, xruntime.Location())
	}
	if p.Base.RunMode == nil {
		return errors.WithMessage(err, xruntime.Location())
	}
	if p.Base.AvailableLoad == nil {
		return errors.WithMessage(err, xruntime.Location())
	}
	if p.Timer.ScanSecondDuration == nil {
		defaultValue := xtimer.ScanSecondDurationDefault
		p.Timer.ScanSecondDuration = &defaultValue
	}
	if p.Timer.ScanMillisecondDuration == nil {
		defaultValue := xtimer.ScanMillisecondDurationDefault
		p.Timer.ScanMillisecondDuration = &defaultValue
	}
	for _, v := range p.ServerNet {
		if v.Addr == nil {
			defaultValue := ""
			v.Addr = &defaultValue
		}
		if v.Type == nil {
			defaultValue := xnetcommon.ServerNetTypeNameTCP
			v.Type = &defaultValue
		}
		if *v.Type != xnetcommon.ServerNetTypeNameTCP && *v.Type != xnetcommon.ServerNetTypeNameKCP {
			return xerror.NotImplemented.WithExtraMessage(fmt.Sprintf("serviceNet.type must be tcp or kcp. %x", xruntime.Location()))
		}
	}
	return nil
}

type Base struct {
	ProjectName         *string `yaml:"projectName"`         // 项目名称
	Version             *string `yaml:"version"`             // 版本号
	PprofHttpPort       *uint16 `yaml:"pprofHttpPort"`       // pprof性能分析 http端口 [default]: nil 不使用
	LogLevel            *uint32 `yaml:"logLevel"`            // 日志等级
	LogAbsPath          *string `yaml:"logAbsPath"`          // 日志绝对路径 [default]: 当前执行的程序-绝对路径,指向启动当前进程的可执行文件-目录路径. e.g.:absPath/log
	GoMaxProcess        *int    `yaml:"goMaxProcess"`        // [default]: runtime.NumCPU()
	BusChannelCapacity  *uint32 `yaml:"busChannelCapacity"`  // 总线chan容量
	PacketLengthMax     *uint32 `yaml:"packetLengthMax"`     // bytes,用户 上行 每个包的最大长度
	SendChannelCapacity *uint32 `yaml:"sendChannelCapacity"` // bytes,每个TCP链接的发送chan大小
	RunMode             *uint32 `yaml:"runMode"`             // 运行模式 [0:release 1:debug]
	AvailableLoad       *uint32 `yaml:"availableLoad"`       // 剩余可用负载, 可用资源数
}

type Timer struct {
	// 秒级定时器 扫描间隔(纳秒) 1000*1000*100=100000000 为100毫秒 [default]: xtimer.ScanSecondDurationDefault
	ScanSecondDuration *time.Duration `yaml:"scanSecondDuration"`
	// 毫秒级定时器 扫描间隔(纳秒) 1000*1000*100=100000000 为25毫秒 [default]: xtimer.ScanMillisecondDurationDefault
	ScanMillisecondDuration *time.Duration `yaml:"scanMillisecondDuration"`
}
