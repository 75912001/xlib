package log

import (
	"context"
	"time"
)

//日志条目

// ExtendFields 扩展字段,日志数据字段
type ExtendFields []interface{} // key,val 数组

// entry 日志数据信息
type entry struct {
	level        uint32    // 本条目的日志级别
	time         time.Time // 生成日志的时间
	callerInfo   string    // 调用堆栈信息
	message      string    // 日志信息
	ctx          context.Context
	extendFields ExtendFields // [string,interface{}] key,value;key,value...
	outString    string       // 输出的字符串
}

func newEntry() *entry {
	return &entry{}
}

func (p *entry) reset() {
	p.level = LevelOff
	p.callerInfo = ""
	p.message = ""
	p.ctx = nil
	p.extendFields = nil
	p.outString = ""
}

func (p *entry) withLevel(level uint32) *entry {
	p.level = level
	return p
}

func (p *entry) withTime(nowTime time.Time) *entry {
	p.time = nowTime
	return p
}

func (p *entry) withCallerInfo(callerInfo string) *entry {
	p.callerInfo = callerInfo
	return p
}

func (p *entry) withMessage(message string) *entry {
	p.message = message
	return p
}

func (p *entry) WithContext(ctx context.Context) *entry {
	p.ctx = ctx
	return p
}

//func (p *entry) WithExtendField(key string, value interface{}) *entry {
//	if p.ExtendFields == nil {
//		p.ExtendFields = make(ExtendFields, 0, 4)
//	}
//	p.ExtendFields = append(p.ExtendFields, key, value)
//	return p
//}

func (p *entry) WithExtendFields(fields ExtendFields) *entry {
	if p.extendFields == nil {
		fieldsSize := len(fields)
		p.extendFields = make(ExtendFields, 0, fieldsSize)
	}
	p.extendFields = append(p.extendFields, fields...)
	return p
}
