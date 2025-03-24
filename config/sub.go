package config

import (
	xruntime "github.com/75912001/xlib/runtime"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
	"time"
)

type Sub struct {
	*Jaeger  `yaml:"jaeger"`  // todo
	*MongoDB `yaml:"mongoDB"` // todo
	*Redis   `yaml:"redis"`   // todo
	*NATS    `yaml:"nats"`    // todo
}

func (p *Sub) Unmarshal(strYaml string) error {
	if err := yaml.Unmarshal([]byte(strYaml), &p); err != nil {
		return errors.WithMessage(err, xruntime.Location())
	}
	return nil
}

func (p *Sub) GetJaeger() *Jaeger {
	return p.Jaeger
}

func (p *Sub) GetMongoDB() *MongoDB {
	return p.MongoDB
}

func (p *Sub) GetRedis() *Redis {
	return p.Redis
}

func (p *Sub) GetNATS() *NATS {
	return p.NATS
}

type Jaeger struct {
	Addrs []string `yaml:"addrs"`
}

type MongoDB struct {
	Addrs           []string       `yaml:"addrs"`
	User            *string        `yaml:"user"`
	Password        *string        `yaml:"password"`
	DBName          *string        `yaml:"dbName"`          // 数据库名称 [default]: todo
	MaxPoolSize     *uint64        `yaml:"maxPoolSize"`     // 连接池最大数量,该数量应该与并发数量匹配 [default]: todo
	MinPoolSize     *uint64        `yaml:"minPoolSize"`     // 池最小数量 [default]: todo
	TimeoutDuration *time.Duration `yaml:"timeoutDuration"` // 操作超时时间 [default]: todo
	MaxConnIdleTime *time.Duration `yaml:"maxConnIdleTime"` // 指定连接在连接池中保持空闲的最长时间 [default]: todo
	MaxConnecting   *uint64        `yaml:"maxConnecting"`   // 指定连接池可以同时建立的最大连接数 [default]: todo
	DBAsync         *DBAsync       `yaml:"dbAsync"`         // DB异步消费配置
}

type Redis struct {
	Addrs    []string `yaml:"addrs"`
	Password *string  `yaml:"password"`
}

type NATS struct {
	Addrs    []string `yaml:"addrs"`
	User     *string  `yaml:"user"`     // 用户 default: todo
	Password *string  `yaml:"password"` // 密码 default: todo
}

type DBAsync struct {
	ChanCnt              *uint32 `yaml:"chanCnt"`              // DB异步消费chan数量. 为0或者没有则不开启异步消费
	Model                *uint32 `yaml:"model"`                // DB异步消费模型 [default] todo
	BulkWriteMax         *uint32 `yaml:"bulkWriteMax"`         // DB合并写 单个集合最大合批数量 [default] todo
	BulkWriteMillisecond *uint32 `yaml:"bulkWriteMillisecond"` // DB合并写周期  单位毫秒 [default] todo
}
