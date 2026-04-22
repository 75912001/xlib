package resolve

import (
	"sync"

	xgrpcutil "github.com/75912001/xlib/grpc/util"
	xhashring "github.com/75912001/xlib/hashring"
	xmap "github.com/75912001/xlib/map"
)

var gHashRingMgr = newHashRingMgr()

type hashRing struct {
	// writer 串行化 add/del 中的「Find → 构造新环 → Add」，避免仅替换 Map 时并发丢更新
	writer sync.Mutex
	m      *xmap.MapMutexMgr[string, *xhashring.HashRing[string]] // key: genPackageNameServiceName() val: 哈希环(node:ServerKey.String())
}

func newHashRingMgr() *hashRing {
	return &hashRing{
		m: xmap.NewMapMutexMgr[string, *xhashring.HashRing[string]](),
	}
}

func (p *hashRing) add(packageName string, serviceName string, node string) {
	packageNameServiceName := genPackageNameServiceName(packageName, serviceName)
	p.writer.Lock()
	defer p.writer.Unlock()

	oldHashRing, ok := p.m.Find(packageNameServiceName)
	if !ok {
		newHashRing := xhashring.NewHashRing[string]().AddNode(node)
		p.m.Add(packageNameServiceName, newHashRing)
		return
	}
	newHashRing := oldHashRing.AddNode(node)
	p.m.Add(packageNameServiceName, newHashRing)
}

func (p *hashRing) del(packageName string, serviceName string, node string) {
	packageNameServiceName := genPackageNameServiceName(packageName, serviceName)
	p.writer.Lock()
	defer p.writer.Unlock()

	v, ok := p.m.Find(packageNameServiceName)
	if !ok {
		return
	}
	newHashRing := v.RemoveNode(node)
	p.m.Add(packageNameServiceName, newHashRing)
}

func (p *hashRing) get(packageName string, serviceName string, shareKey string) xgrpcutil.IClientConn {
	packageNameServiceName := genPackageNameServiceName(packageName, serviceName)
	v, ok := p.m.Find(packageNameServiceName)
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
