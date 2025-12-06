// 服务配置
// 服务配置文件, 用于配置服务的基本信息.
// 该配置文件与可执行程序在同一目录下.

package config

import (
	xruntime "github.com/75912001/xlib/runtime"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
	"os"
)

var GConfigMgr = NewMgr()

// 配置-主项,用户服务的基本配置

type Mgr struct {
	ExecutablePath string         // 绝对路径
	Content        string         // 内容
	Etcd           Etcd           `yaml:"etcd"`
	Base           Base           `yaml:"base"`
	Log            Log            `yaml:"log"`
	Timer          Timer          `yaml:"timer"`
	Net            []*Net         `yaml:"net"`
	KCP            KCP            `yaml:"kcp"`
	Grpc           Grpc           `yaml:"grpc"`
	Redis          []*Redis       `yaml:"redis"`
	Nats           []*Nats        `yaml:"nats"`
	Custom         map[string]any `yaml:"custom"` // 自定义配置，支持各模块自行解析
	//# 自定义配置区域
	//custom:
	//# string
	//testString: hello
	//# int/uint32
	//maxConnections: 1000
	//# bool
	//enableDebug: true
}

func NewMgr() *Mgr {
	return &Mgr{}
}

func (p *Mgr) Parse(executablePath string) error {
	p.ExecutablePath = executablePath
	content, err := os.ReadFile(p.ExecutablePath)
	if err != nil {
		return errors.WithMessagef(err, "read file %v failed. %v", p.ExecutablePath, xruntime.Location())
	}
	p.Content = string(content)

	if err := yaml.Unmarshal([]byte(p.Content), &p); err != nil {
		return errors.WithMessagef(err, "unmarshal file %v failed. %v", p.ExecutablePath, xruntime.Location())
	}
	if err := p.Etcd.Configure(); err != nil {
		return errors.WithMessagef(err, "etcd configure failed. %v", xruntime.Location())
	}
	if err := p.Base.Configure(); err != nil {
		return errors.WithMessagef(err, "base configure failed. %v", xruntime.Location())
	}
	if err := p.Log.Configure(); err != nil {
		return errors.WithMessagef(err, "log configure failed. %v", xruntime.Location())
	}
	if err := p.Timer.Configure(); err != nil {
		return errors.WithMessagef(err, "timer configure failed. %v", xruntime.Location())
	}
	for _, v := range p.Net {
		if err := v.Configure(); err != nil {
			return errors.WithMessagef(err, "net configure failed. %v", xruntime.Location())
		}
	}
	if err := p.KCP.Configure(); err != nil {
		return errors.WithMessagef(err, "kcp configure failed. %v", xruntime.Location())
	}
	if err := p.Grpc.Configure(); err != nil {
		return errors.WithMessagef(err, "grpc configure failed. %v", xruntime.Location())
	}
	for _, v := range p.Redis {
		if err := v.Configure(); err != nil {
			return errors.WithMessagef(err, "redis configure failed. %v", xruntime.Location())
		}
	}
	for _, v := range p.Nats {
		if err := v.Configure(); err != nil {
			return errors.WithMessagef(err, "nats configure failed. %v", xruntime.Location())
		}
	}
	return nil
}

// GetCustomUint32 获取 uint32 类型的配置值
func (p *Mgr) GetCustomUint32(key string, defaultValue ...uint32) uint32 {
	var dv uint32
	if len(defaultValue) == 0 {
		dv = 0
	} else {
		dv = defaultValue[0]
	}
	if p.Custom == nil {
		return dv
	}

	val, exists := p.Custom[key]
	if !exists {
		return dv
	}

	// 处理多种可能的类型
	switch v := val.(type) {
	case int:
		return uint32(v)
	case int32:
		return uint32(v)
	case int64:
		return uint32(v)
	case uint32:
		return v
	case uint64:
		return uint32(v)
	case float64: // YAML 解析数字时可能是 float64
		return uint32(v)
	default:
		return dv
	}
}

func (p *Mgr) GetCustomInt64(key string, defaultValue ...int64) int64 {
	var dv int64
	if len(defaultValue) == 0 {
		dv = 0
	} else {
		dv = defaultValue[0]
	}
	if p.Custom == nil {
		return dv
	}

	val, exists := p.Custom[key]
	if !exists {
		return dv
	}

	switch v := val.(type) {
	case int:
		return int64(v)
	case int32:
		return int64(v)
	case int64:
		return v
	case uint32:
		return int64(v)
	case uint64:
		return int64(v)
	case float64:
		return int64(v)
	default:
		return dv
	}
}

func (p *Mgr) GetCustomUint64(key string, defaultValue ...uint64) uint64 {
	var dv uint64
	if len(defaultValue) == 0 {
		dv = 0
	} else {
		dv = defaultValue[0]
	}
	if p.Custom == nil {
		return dv
	}

	val, exists := p.Custom[key]
	if !exists {
		return dv
	}

	// 处理多种可能的类型
	switch v := val.(type) {
	case int:
		return uint64(v)
	case int32:
		return uint64(v)
	case int64:
		return uint64(v)
	case uint32:
		return uint64(v)
	case uint64:
		return v
	case float64: // YAML 解析数字时可能是 float64
		return uint64(v)
	default:
		return dv
	}
}

// GetCustomInt 获取 int 类型的配置值（额外提供）
func (p *Mgr) GetCustomInt(key string, defaultValue ...int) int {
	var dv int
	if len(defaultValue) == 0 {
		dv = 0
	} else {
		dv = defaultValue[0]
	}
	if p.Custom == nil {
		return dv
	}

	val, exists := p.Custom[key]
	if !exists {
		return dv
	}

	switch v := val.(type) {
	case int:
		return v
	case int32:
		return int(v)
	case int64:
		return int(v)
	case uint32:
		return int(v)
	case uint64:
		return int(v)
	case float64:
		return int(v)
	default:
		return dv
	}
}

// GetCustomString 获取 string 类型的配置值
func (p *Mgr) GetCustomString(key string, defaultValue ...string) string {
	var dv string
	if len(defaultValue) == 0 {
		dv = ""
	} else {
		dv = defaultValue[0]
	}
	if p.Custom == nil {
		return dv
	}

	val, exists := p.Custom[key]
	if !exists {
		return dv
	}

	if str, ok := val.(string); ok {
		return str
	}

	return dv
}

// GetCustomBool 获取 bool 类型的配置值（额外提供）
func (p *Mgr) GetCustomBool(key string, defaultValue ...bool) bool {
	var dv bool
	if len(defaultValue) == 0 {
		dv = false
	} else {
		dv = defaultValue[0]
	}
	if p.Custom == nil {
		return dv
	}

	val, exists := p.Custom[key]
	if !exists {
		return dv
	}

	if b, ok := val.(bool); ok {
		return b
	}

	return dv
}
