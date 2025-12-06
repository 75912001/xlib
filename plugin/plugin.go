package plugin

import (
	xerror "github.com/75912001/xlib/error"
	xmap "github.com/75912001/xlib/map"
	xruntime "github.com/75912001/xlib/runtime"
	"github.com/pkg/errors"
	"sync"
)

// IPlugin 配置插件接口
type IPlugin interface {
	// Name 返回插件名称
	Name() string
	// Init 初始化插件
	Init() error
	// Close 关闭插件
	Close() error
}

// Mgr 插件管理器
type Mgr struct {
	pluginMapMgr *xmap.MapMgr[string, IPlugin] // key: 插件名称, value: 插件实例
	mu           sync.RWMutex
}

// NewPluginManager 创建插件管理器
func NewPluginManager() *Mgr {
	return &Mgr{
		pluginMapMgr: xmap.NewMapMgr[string, IPlugin](),
	}
}

// Register 注册插件
func (p *Mgr) Register(plugin IPlugin) error {
	if plugin == nil {
		return errors.WithMessagef(xerror.Param, "plugin is nil. %v", xruntime.Location())
	}
	name := plugin.Name()
	if name == "" {
		return errors.WithMessagef(xerror.Param, "plugin name cannot be empty. %v", xruntime.Location())
	}

	p.mu.Lock()
	defer func() {
		p.mu.Unlock()
	}()

	if _, exists := p.pluginMapMgr.Find(name); exists {
		return errors.WithMessagef(xerror.Param, "plugin %v already registered. %v", name, xruntime.Location())
	}

	p.pluginMapMgr.Add(name, plugin)
	return nil
}

// Get 获取插件
func (p *Mgr) Get(name string) IPlugin {
	p.mu.RLock()
	defer func() {
		p.mu.RUnlock()
	}()

	return p.pluginMapMgr.Get(name)
}

// Remove 移除插件-关闭插件并从管理器中删除
func (p *Mgr) Remove(name string) error {
	p.mu.Lock()
	defer func() {
		p.mu.Unlock()
	}()

	if plugin, exists := p.pluginMapMgr.Find(name); exists { // 有
		if err := plugin.Close(); err != nil {
			return errors.WithMessagef(err, "close plugin %v error. %v", name, xruntime.Location())
		}
		p.pluginMapMgr.Del(name)
	}
	return nil
}

// Close 关闭所有插件
func (p *Mgr) Close() error {
	var returnError error

	p.mu.RLock()
	defer func() {
		p.mu.RUnlock()
	}()

	p.pluginMapMgr.Foreach(
		func(name string, plugin IPlugin) bool {
			if err := plugin.Close(); err != nil {
				if returnError == nil {
					returnError = errors.WithMessagef(err, "close plugin %v error.", name)
				} else {
					returnError = errors.WithMessagef(returnError, "%v close plugin %v error.", err, name)
				}
			}
			return true
		},
	)
	if returnError != nil {
		returnError = errors.WithMessage(returnError, xruntime.Location())
	}
	return returnError
}

// List 列出所有已注册的插件
func (p *Mgr) List() []string {
	var names []string

	p.mu.RLock()
	defer func() {
		p.mu.RUnlock()
	}()

	p.pluginMapMgr.Foreach(
		func(name string, plugin IPlugin) bool {
			names = append(names, name)
			return true
		})
	return names
}

// Count 返回已注册的插件数量
func (p *Mgr) Count() int {
	p.mu.RLock()
	defer func() {
		p.mu.RUnlock()
	}()

	return p.pluginMapMgr.Len()
}
