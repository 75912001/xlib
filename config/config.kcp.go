package config

type KCP struct {
	Password *string `yaml:"password"` // 密码		[default]: "demo.pass"
	Salt     *string `yaml:"salt"`     // 盐		[default]: "demo.salt"
}

func (p *KCP) Configure() error {
	if p.Password == nil {
		defaultValue := "demo.pass"
		p.Password = &defaultValue
	}
	if p.Salt == nil {
		defaultValue := "demo.salt"
		p.Salt = &defaultValue
	}
	return nil
}
