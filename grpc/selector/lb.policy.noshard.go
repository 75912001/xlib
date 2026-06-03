package selector

import (
	"context"
	"sync"
	"sync/atomic"

	xerror "github.com/75912001/xlib/error"
	xgrpcresolve "github.com/75912001/xlib/grpc/resolve"
	"github.com/75912001/xlib/grpc/util"
	xgrpcutil "github.com/75912001/xlib/grpc/util"
	xruntime "github.com/75912001/xlib/runtime"
	xutil "github.com/75912001/xlib/util"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

// First 选择当前服务列表中的第一个连接。
// 该策略不依赖 shardKey，但仍由 selector 从服务发现连接列表中选择目标连接。
type First[K util.IKey] struct {
}

func newFirst[K util.IKey]() *First[K] {
	return &First[K]{}
}

// Select 返回当前方法所属服务的第一个可用连接。
func (p *First[K]) Select(ctx context.Context, key K, method string) (*grpc.ClientConn, error) {
	clientConnSlice, err := getClientConnSlice(method)
	if err != nil {
		return nil, err
	}
	return clientConnSlice[0].GetClientConn(), nil
}

// Random 从当前服务列表中随机选择一个连接。
// 该策略不依赖 shardKey，适合无状态或不要求同 key 固定落点的 unary RPC。
type Random[K util.IKey] struct {
}

func newRandom[K util.IKey]() *Random[K] {
	return &Random[K]{}
}

// Select 随机返回当前方法所属服务的一个连接。
func (p *Random[K]) Select(ctx context.Context, key K, method string) (*grpc.ClientConn, error) {
	clientConnSlice, err := getClientConnSlice(method)
	if err != nil {
		return nil, err
	}
	idx := int(xutil.SecureRandomUint64() % uint64(len(clientConnSlice)))
	return clientConnSlice[idx].GetClientConn(), nil
}

// RoundRobin 按方法维度轮询选择服务连接。
// 每个 method 独立维护计数器，避免不同 RPC 方法共享同一个轮询位置。
type RoundRobin[K util.IKey] struct {
	counters sync.Map
}

func newRoundRobin[K util.IKey]() *RoundRobin[K] {
	return &RoundRobin[K]{}
}

// Select 轮询返回当前方法所属服务的一个连接。
func (p *RoundRobin[K]) Select(ctx context.Context, key K, method string) (*grpc.ClientConn, error) {
	clientConnSlice, err := getClientConnSlice(method)
	if err != nil {
		return nil, err
	}
	counter := p.getCounter(method)
	idx := int((counter.Add(1) - 1) % uint64(len(clientConnSlice)))
	return clientConnSlice[idx].GetClientConn(), nil
}

func (p *RoundRobin[K]) getCounter(method string) *atomic.Uint64 {
	counter, _ := p.counters.LoadOrStore(method, &atomic.Uint64{})
	return counter.(*atomic.Uint64)
}

// getClientConnSlice 解析完整 RPC method，并从服务发现缓存中获取对应服务的连接列表。
func getClientConnSlice(method string) ([]xgrpcutil.IClientConn, error) {
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
	packetServiceName := "/" + m.PackageName + "." + m.ServiceName
	clientConnSlice := xgrpcresolve.GetClientConn(packetServiceName)
	if len(clientConnSlice) == 0 {
		return nil, errors.WithMessage(xerror.NotFound, xruntime.Location())
	}
	return clientConnSlice, nil
}
