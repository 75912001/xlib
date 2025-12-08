package module

import (
	xruntime "github.com/75912001/xlib/runtime"
	"github.com/pkg/errors"
)

// Factory 模块工厂结构体
//
//	采用数组的方式, 索引:模块类型, 值:模块元素的指针
type Factory struct {
	elements []*Module
}

// NewFactory 创建新的工厂实例
//
//	参数: capacity - 模块创建器数组的初始容量
//	e.g.: 0: a模块, 1: b模块, 2: c模块 => 传递 capacity = 3
func NewFactory(capacity uint32) *Factory {
	return &Factory{
		elements: make([]*Module, capacity),
	}
}

// Register 注册模块
func (p *Factory) Register(element *Module) error {
	if uint32(len(p.elements)) <= uint32(element.Type) {
		return errors.WithMessagef(errors.New("module type out of range"), "module type: %v %v", element.Type, xruntime.Location())
	}
	if p.elements[element.Type] != nil {
		return errors.WithMessagef(errors.New("module type already registered"), "module type: %v %v", element.Type, xruntime.Location())
	}
	if element.CreatorFunction == nil {
		return errors.WithMessagef(errors.New("module creator function is nil"), "module type: %v %v", element.Type, xruntime.Location())
	}
	p.elements[element.Type] = element
	return nil
}

// CreateAllModule 创建所有模块
func (p *Factory) CreateAllModule(holder IHolder) {
	for i := 0; i < len(p.elements); i++ {
		element := p.elements[i]
		if element == nil {
			continue
		}
		element.IModule = element.CreatorFunction(holder)
	}
}

// GetModule 获取模块
func (p *Factory) GetModule(moduleType Type) *Module {
	return p.elements[moduleType]
}

// Foreach 遍历所有模块
func (p *Factory) Foreach(callback func(element *Module)) {
	for _, element := range p.elements {
		if element == nil {
			continue
		}
		callback(element)
	}
}
