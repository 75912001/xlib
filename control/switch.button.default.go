package control

type SwitchButton struct {
	state bool // 状态 true:on false:off
}

func NewSwitchButton(state bool) *SwitchButton {
	return &SwitchButton{
		state: state,
	}
}

func (p *SwitchButton) On() {
	p.state = true
}

func (p *SwitchButton) Off() {
	p.state = false
}

func (p *SwitchButton) IsOn() bool {
	return p.state
}

func (p *SwitchButton) IsOff() bool {
	return !p.state
}
