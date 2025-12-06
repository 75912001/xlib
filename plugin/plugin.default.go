package plugin

type Default struct {
	name string
}

func NewDefault(name string) *Default {
	return &Default{
		name: name,
	}
}

func (p *Default) Name() string {
	return p.name
}

func (p *Default) Init() error {
	return nil
}

func (p *Default) Close() error {
	return nil
}
