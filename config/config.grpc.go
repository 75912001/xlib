package config

import (
	"math"

	xerror "github.com/75912001/xlib/error"
	xruntime "github.com/75912001/xlib/runtime"
	"github.com/pkg/errors"
)

const (
	grpcMaxReceiveMessageBytesDefault = 4 * 1024 * 1024
	grpcMaxSendMessageBytesDefault    = math.MaxInt32
)

// Grpc gRPC 服务配置
// 通过
type Grpc struct {
	PackageName            *string `yaml:"packageName"`            // 包名
	ServiceName            *string `yaml:"serviceName"`            // 服务名称
	ListenAddr             *string `yaml:"listenAddr"`             // 服务地址-Listen (如果配置,则Listen服务) e.g.: 127.0.0.1:8989		[default]: "127.0.0.1:6523"
	ExternalAddr           *string `yaml:"externalAddr"`           // 服务地址-对外 e.g.: 127.0.0.1:8989		[default]: 未配置-使用 -> 服务地址-Listen
	MaxReceiveMessageBytes *int    `yaml:"maxReceiveMessageBytes"` // 单条 gRPC 接收消息上限, bytes
	MaxSendMessageBytes    *int    `yaml:"maxSendMessageBytes"`    // 单条 gRPC 发送消息上限, bytes
}

func (p *Grpc) GetMaxReceiveMessageBytes() int {
	if p.MaxReceiveMessageBytes != nil {
		return *p.MaxReceiveMessageBytes
	}
	return grpcMaxReceiveMessageBytesDefault
}

func (p *Grpc) GetMaxSendMessageBytes() int {
	if p.MaxSendMessageBytes != nil {
		return *p.MaxSendMessageBytes
	}
	return grpcMaxSendMessageBytesDefault
}

func (p *Grpc) HasListenAddr() bool {
	return p.ListenAddr != nil && *p.ListenAddr != ""
}

// 是否启用
func (p *Grpc) IsEnabled() bool {
	return p.HasListenAddr()
}

func (p *Grpc) Configure() error {
	if p.PackageName == nil {
		defaultValue := ""
		p.PackageName = &defaultValue
	}
	if p.ServiceName == nil {
		defaultValue := ""
		p.ServiceName = &defaultValue
	}
	if p.ListenAddr == nil {
		defaultValue := ""
		p.ListenAddr = &defaultValue
	} else {
		if *p.PackageName == "" || *p.ServiceName == "" { // 地址已配置, 缺少包名或服务名
			return errors.WithMessagef(xerror.Configure, "packageName or serviceName is empty. %v", xruntime.Location())
		}
	}
	if p.ExternalAddr == nil {
		p.ExternalAddr = p.ListenAddr
	}
	if p.MaxReceiveMessageBytes == nil {
		defaultValue := grpcMaxReceiveMessageBytesDefault
		p.MaxReceiveMessageBytes = &defaultValue
	} else if *p.MaxReceiveMessageBytes <= 0 {
		return errors.WithMessagef(xerror.Configure, "maxReceiveMessageBytes must be greater than 0. %v", xruntime.Location())
	}
	if p.MaxSendMessageBytes == nil {
		defaultValue := grpcMaxSendMessageBytesDefault
		p.MaxSendMessageBytes = &defaultValue
	} else if *p.MaxSendMessageBytes <= 0 {
		return errors.WithMessagef(xerror.Configure, "maxSendMessageBytes must be greater than 0. %v", xruntime.Location())
	}
	return nil
}
