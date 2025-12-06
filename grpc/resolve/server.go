package resolve

import (
	"fmt"
	xgrpcutil "github.com/75912001/xlib/grpc/util"
	xmap "github.com/75912001/xlib/map"
)

var GServerMgr = newServerMgr()

type ServerKey struct {
	GroupID    uint32
	ServerName string
	ServerID   uint32
}

func (p *ServerKey) String() string {
	// uint32 转为字符串时, 不够9位, 需要补0
	return fmt.Sprintf("%09d.%v.%09d", p.GroupID, p.ServerName, p.ServerID)
}

type ServerMgr struct {
	*xmap.MapMutexMgr[string, xgrpcutil.IClientConn] // 存储 server 对应的 gRPC 客户端连接 key: ServerKey.String() val: gRPC 客户端连接
}

func newServerMgr() *ServerMgr {
	return &ServerMgr{
		MapMutexMgr: xmap.NewMapMutexMgr[string, xgrpcutil.IClientConn](),
	}
}
