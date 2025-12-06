package resources

var GResources = NewResources()

type Resources struct {
	availableLoad uint32 // 可用负载
}

// NewResources 新的 Resources
func NewResources() *Resources {
	return &Resources{}
}

// 获取可用负载
func (p *Resources) GetAvailableLoad() uint32 {
	return p.availableLoad
}

// 设置可用负载
func (p *Resources) SetAvailableLoad(availableLoad uint32) {
	p.availableLoad = availableLoad
}

// 增加可用负载
func (p *Resources) AddAvailableLoad(delta uint32) {
	p.availableLoad += delta
}

// 减少可用负载
func (p *Resources) SubAvailableLoad(delta uint32) {
	if p.availableLoad < delta {
		p.availableLoad = 0
	} else {
		p.availableLoad -= delta
	}
}
