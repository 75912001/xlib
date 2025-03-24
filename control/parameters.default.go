package control

type Parameters struct {
	parameters []interface{}
}

func NewParameters() *Parameters {
	return &Parameters{
		parameters: make([]interface{}, 0),
	}
}

func (p *Parameters) Override(parameters ...interface{}) {
	p.parameters = append([]interface{}{}, parameters...)
}

func (p *Parameters) Get() []interface{} {
	return p.parameters
}

func (p *Parameters) Append(parameters ...interface{}) {
	p.parameters = append(p.parameters, parameters...)
}
