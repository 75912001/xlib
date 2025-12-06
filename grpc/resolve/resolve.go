package resolve

import (
	xerror "github.com/75912001/xlib/error"
	xgrpcutil "github.com/75912001/xlib/grpc/util"
	xruntime "github.com/75912001/xlib/runtime"
	"github.com/pkg/errors"
)

// 服务失效, 删除服务
func RemoveServer(groupID uint32, serverName string, serverID uint32, packageName string, serviceName string) (xgrpcutil.IClientConn, error) {
	serverKey := ServerKey{
		GroupID:    groupID,
		ServerName: serverName,
		ServerID:   serverID,
	}
	conn, ok := GServerMgr.Find(serverKey.String())
	if !ok {
		return nil, errors.WithMessage(xerror.GRPCServiceNotFound, xruntime.Location())
	}
	GServerMgr.Del(serverKey.String())
	conn.Disabled()
	_ = conn.Stop()
	gPacketServiceMgr.del(conn)
	gHashRingMgr.del(packageName, serviceName, serverKey.String())
	return conn, nil
}

// 服务上线, 添加服务
func AddServer(groupID uint32, serverName string, serverID uint32, conn xgrpcutil.IClientConn, packageName string, serviceName string) {
	serverKey := ServerKey{
		GroupID:    groupID,
		ServerName: serverName,
		ServerID:   serverID,
	}
	GServerMgr.Add(serverKey.String(), conn)
	gPacketServiceMgr.add(packageName, serviceName, conn)
	gHashRingMgr.add(packageName, serviceName, serverKey.String())
}

// 获取服务连接
func GetClientConn(packetServiceName string) []xgrpcutil.IClientConn {
	return gPacketServiceMgr.get(packetServiceName)
}

func GetClientConnByHashRing(packageName string, serviceName string, shardKey string) (xgrpcutil.IClientConn, error) {
	conn := gHashRingMgr.get(packageName, serviceName, shardKey)
	if conn == nil {
		return nil, errors.WithMessage(xerror.NotExist, xruntime.Location())
	}
	return conn, nil
}
