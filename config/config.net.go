package config

import (
	xerror "github.com/75912001/xlib/error"
	xnetcommon "github.com/75912001/xlib/net/common"
	xruntime "github.com/75912001/xlib/runtime"
	"github.com/pkg/errors"
)

type Net struct {
	Name         *string `yaml:"name"`         // 链接名称		[default]: "netName"
	Type         *string `yaml:"type"`         // [tcp, kcp, websocket]		[default]: common.ServerNetTypeNameTCP
	ListenAddr   *string `yaml:"listenAddr"`   // 服务地址-Listen (如果配置,则Listen服务) e.g.: 127.0.0.1:8989
	ExternalAddr *string `yaml:"externalAddr"` // 服务地址-对外 e.g.: 127.0.0.1:8989		[default]: 未配置-使用 -> 服务地址-Listen
	Pattern      *string `yaml:"pattern"`      // 用于 type: websocket
}

func (p *Net) Configure() error {
	if p.Name == nil {
		defaultValue := "netName"
		p.Name = &defaultValue
	}
	if p.Type == nil {
		defaultValue := xnetcommon.ServerNetTypeNameTCP
		p.Type = &defaultValue
	}
	if *p.Type != xnetcommon.ServerNetTypeNameTCP &&
		*p.Type != xnetcommon.ServerNetTypeNameKCP &&
		*p.Type != xnetcommon.ServerNetTypeNameWebSocket {
		return errors.WithMessagef(xerror.NotImplemented, "serviceNet.type must be tcp || kcp || websocket. %v", xruntime.Location())
	}
	if p.ListenAddr == nil {
		return errors.WithMessagef(xerror.Config, "serviceNet.listenAddr is empty.")
	}
	if p.ExternalAddr == nil {
		p.ExternalAddr = p.ListenAddr
	}
	switch *p.Type {
	case xnetcommon.ServerNetTypeNameWebSocket:
		if p.Pattern == nil {
			return errors.WithMessagef(xerror.Configure, "net websocket pattern must be set. %v", xruntime.Location())
		}
	}
	return nil
}
