package selector

import (
	"context"
	xerror "github.com/75912001/xlib/error"
	xgrpcresolve "github.com/75912001/xlib/grpc/resolve"
	"github.com/75912001/xlib/grpc/util"
	xgrpcutil "github.com/75912001/xlib/grpc/util"
	xhashring "github.com/75912001/xlib/hashring"
	xruntime "github.com/75912001/xlib/runtime"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"strconv"
)

type HashRing[K util.IKey] struct {
	mgr *xhashring.HashRing[string]
}

func newHashRing[K util.IKey]() *HashRing[K] {
	return &HashRing[K]{
		mgr: xhashring.NewHashRing[string](),
	}
}

func (p *HashRing[K]) Select(ctx context.Context, key K, method string) (*grpc.ClientConn, error) {
	m := xgrpcutil.Method{
		Method:      method,
		PackageName: "",
		ServiceName: "",
		MethodName:  "",
	}
	err := m.Parse()
	if err != nil {
		return nil, errors.WithMessagef(xerror.GRPCInvalidMethod, "method %s parse error: %v", method, err)
	}
	var strKey string
	switch v := any(key).(type) {
	case string:
		strKey = v
	case int32:
		strKey = strconv.FormatInt(int64(v), 10)
	case int64:
		strKey = strconv.FormatInt(v, 10)
	case uint32:
		strKey = strconv.FormatUint(uint64(v), 10)
	case uint64:
		strKey = strconv.FormatUint(v, 10)
	default:
		return nil, errors.WithMessage(xerror.NotSupport, xruntime.Location())
	}
	clientConn, err := xgrpcresolve.GetClientConnByHashRing(m.PackageName, m.ServiceName, strKey)
	if err != nil {
		return nil, errors.WithMessagef(xerror.NotExist, "err: %v", err)
	}
	return clientConn.GetClientConn(), nil
}
