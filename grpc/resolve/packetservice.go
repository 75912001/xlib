package resolve

import (
	xgrpcutil "github.com/75912001/xlib/grpc/util"
	xmap "github.com/75912001/xlib/map"
	"sort"
	"sync"
)

// "/packageName.serviceName"

var gPacketServiceMgr = newPacketServiceMgr()

type packetServiceMgr struct {
	mapMgr *xmap.MapMgr[string, []xgrpcutil.IClientConn] // key: "/packageName.serviceName" val: 连接列表(连接列表中的连接, 需要排序,否则,不用的服务中,使用相同的 shardKey ,会导致获取到的连接不一致)
	mu     sync.RWMutex
}

func newPacketServiceMgr() *packetServiceMgr {
	return &packetServiceMgr{
		mapMgr: xmap.NewMapMgr[string, []xgrpcutil.IClientConn](),
	}
}

func (p *packetServiceMgr) add(packageName string, serviceName string, clientConn xgrpcutil.IClientConn) {
	key := xgrpcutil.GenPackageServiceName(packageName, serviceName)
	p.mu.Lock()
	defer p.mu.Unlock()

	clientConnSlice, ok := p.mapMgr.Find(key)
	if !ok {
		clientConnSlice = make([]xgrpcutil.IClientConn, 0)
	}
	clientConnSlice = append(clientConnSlice, clientConn)
	// 排序 使用相同的 shardKey ,获取到的连接是相同的
	sort.Slice(clientConnSlice, func(i, j int) bool {
		return clientConnSlice[i].GetID() < clientConnSlice[j].GetID()
	})
	p.mapMgr.Add(key, clientConnSlice)
}

func (p *packetServiceMgr) del(delClientConn xgrpcutil.IClientConn) {
	// 收集需要删除和更新的数据
	toDelete := make([]string, 0)
	toUpdate := make(map[string][]xgrpcutil.IClientConn)

	p.mu.Lock()
	defer p.mu.Unlock()

	p.mapMgr.Foreach(
		func(key string, clientConnSlice []xgrpcutil.IClientConn) (isContinue bool) {
			reservedSlice := make([]xgrpcutil.IClientConn, 0)
			for _, clientConn := range clientConnSlice {
				if delClientConn != clientConn {
					reservedSlice = append(reservedSlice, clientConn)
				}
			}
			if len(reservedSlice) == 0 {
				toDelete = append(toDelete, key)
			} else {
				toUpdate[key] = reservedSlice
			}
			return true
		},
	)

	p.mapMgr.Del(toDelete...)

	for key, clientConnSlice := range toUpdate {
		p.mapMgr.Add(key, clientConnSlice)
	}
}

func (p *packetServiceMgr) get(packetServiceName string) []xgrpcutil.IClientConn {
	p.mu.RLock()
	defer p.mu.RUnlock()

	return p.mapMgr.Get(packetServiceName)
}
