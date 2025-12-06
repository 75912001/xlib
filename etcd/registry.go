package etcd

import (
	xmap "github.com/75912001/xlib/map"
)

var GRegistry = NewRegistry()

// 记录-服务信息
type Registry struct {
	DataMap *xmap.MapMutexMgr[string, *ValueJson]
}

func NewRegistry() *Registry {
	return &Registry{
		DataMap: xmap.NewMapMutexMgr[string, *ValueJson](),
	}
}

func (p *Registry) Update(key string, value *ValueJson) {
	if value == nil { // 如果 value 为空，则删除该 key
		p.DataMap.Del(key)
		return
	}
	p.DataMap.Add(key, value)
}

func (p *Registry) Find(key string) (*ValueJson, bool) {
	return p.DataMap.Find(key)
}

func (p *Registry) FindByGroupNameID(groupID uint32, serviceName string, serviceID uint32) []*ValueJson {
	var results []*ValueJson
	p.DataMap.Foreach(func(key string, value *ValueJson) bool {
		_, _groupID, _serviceName, _serviceID := Parse(key)
		if _groupID == groupID && _serviceName == serviceName && _serviceID == serviceID {
			results = append(results, value)
			return false
		}
		return true
	})
	return results
}

func (p *Registry) FindByGroupName(groupID uint32, serviceName string) []*ServiceInfo {
	var results []*ServiceInfo
	p.DataMap.Foreach(func(key string, value *ValueJson) bool {
		_, _groupID, _serviceName, _serviceID := Parse(key)
		if _groupID == groupID && _serviceName == serviceName {
			results = append(results,
				&ServiceInfo{
					ServiceID:   _serviceID,
					ServiceName: serviceName,
					ValueJson:   value,
				},
			)
		}
		return true
	})
	return results
}

type ServiceInfo struct {
	ServiceID   uint32
	ServiceName string
	ValueJson   *ValueJson
}
