package resolve

import (
	xgrpcutil "github.com/75912001/xlib/grpc/util"
	xhashring "github.com/75912001/xlib/hashring"
	xmap "github.com/75912001/xlib/map"
)

var gHashRingMgr = newHashRingMgr()

type hashRing struct {
	*xmap.MapMgr[string, *xhashring.HashRing[string]] // key: genPackageNameServiceName() val: 哈希环(node:ServerKey.String())
}

func newHashRingMgr() *hashRing {
	return &hashRing{
		MapMgr: xmap.NewMapMgr[string, *xhashring.HashRing[string]](),
	}
}

func (p *hashRing) add(packageName string, serviceName string, node string) {
	packageNameServiceName := genPackageNameServiceName(packageName, serviceName)
	oldHashRing, ok := p.MapMgr.Find(packageNameServiceName)
	if !ok { // 如果不存在，则创建一个新的哈希环
		newHashRing := xhashring.NewHashRing[string]().AddNode(node)
		p.MapMgr.Add(packageNameServiceName, newHashRing)
	} else { // 如果存在，则更新哈希环
		newHashRing := oldHashRing.AddNode(node)
		p.MapMgr.Add(packageNameServiceName, newHashRing)
	}
}

func (p *hashRing) del(packageName string, serviceName string, node string) {
	packageNameServiceName := genPackageNameServiceName(packageName, serviceName)
	v, ok := p.MapMgr.Find(packageNameServiceName)
	if !ok {
		return
	}
	newHashRing := v.RemoveNode(node)
	p.MapMgr.Add(packageNameServiceName, newHashRing)
}

func (p *hashRing) get(packageName string, serviceName string, shareKey string) xgrpcutil.IClientConn {
	packageNameServiceName := genPackageNameServiceName(packageName, serviceName)
	v, ok := p.MapMgr.Find(packageNameServiceName)
	if !ok {
		return nil
	}
	node, ok := v.GetNode(shareKey)
	if !ok {
		return nil
	}
	conn, ok := GServerMgr.Find(node)
	if !ok {
		return nil
	}
	return conn
}
