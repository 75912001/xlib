package control

type Parameters struct {
	args []any
}

func NewParameters() *Parameters {
	return &Parameters{}
}

func (p *Parameters) Override(args ...any) {
	p.args = args
}

func (p *Parameters) Get() []any {
	return p.args
}

func (p *Parameters) Append(args ...any) {
	p.args = append(p.args, args...)
}
