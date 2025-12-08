package util

import (
	xerror "github.com/75912001/xlib/error"
	xruntime "github.com/75912001/xlib/runtime"
	"github.com/pkg/errors"
	"strings"
)

type Method struct {
	Method      string // 完整方法名 "/package.Service/MethodName"
	PackageName string // 包名 "package"
	ServiceName string // 服务名 "Service"
	MethodName  string // 方法名 "MethodName"
}

func NewMethod(method string) *Method {
	return &Method{
		Method: method,
	}
}

func (p *Method) Parse() error {
	if !strings.HasPrefix(p.Method, "/") {
		return errors.WithMessagef(xerror.GRPCInvalidMethod, "method:%v location:%v", p.Method, xruntime.Location())
	}
	parts := strings.Split(p.Method[1:], "/")
	if len(parts) != 2 {
		return errors.WithMessagef(xerror.GRPCInvalidMethod, "method:%v location:%v", p.Method, xruntime.Location())
	}
	pkgService := parts[0]
	dotIdx := strings.LastIndex(pkgService, ".")
	if dotIdx == -1 {
		return errors.WithMessagef(xerror.GRPCInvalidMethod, "method:%v location:%v", p.Method, xruntime.Location())
	}
	p.PackageName = pkgService[:dotIdx]
	p.ServiceName = pkgService[dotIdx+1:]
	p.MethodName = parts[1]
	return nil
}
