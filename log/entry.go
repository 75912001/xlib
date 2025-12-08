package log

import (
	"context"
	"fmt"
	xpool "github.com/75912001/xlib/pool"
	"time"
)

//日志条目

// ExtendFields 扩展字段,日志数据字段
type ExtendFields []any // key,val 数组

// entry 日志数据信息
type entry struct {
	level uint32    // 本条目的日志级别
	time  time.Time // 生成日志的时间

	format string // 日志格式
	args   []any  // 日志参数

	line     int    // 调用堆栈信息-行
	file     string // 调用堆栈信息-文件名
	funcName string // 调用堆栈信息-方法名称

	ctx          context.Context
	extendFields ExtendFields // [string,any] key,value;key,value...
	outString    string       // 输出的字符串
}

func newEntry() *entry {
	return &entry{}
}

func (p *entry) reset() {
	p.level = LevelOff
	p.line = 0
	p.file = ""
	p.funcName = ""
	p.format = ""
	p.args = p.args[0:0]
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

func (p *entry) withCallerInfo(line int, file, funcName string) *entry {
	p.line = line
	p.file = file
	p.funcName = funcName
	return p
}

func (p *entry) getCallerInfo() string {
	return fmt.Sprintf(callerInfoFormat, p.line, p.file, p.funcName)
}

func (p *entry) withMessage(format string, args ...any) *entry {
	p.format = format
	p.args = args
	return p
}

func (p *entry) getMessage() string {
	if p.format != "" {
		return fmt.Sprintf(p.format, p.args...)
	}

	buf := xpool.Buffer.Get()
	defer func() {
		xpool.Buffer.Put(buf)
	}()

	buf.Grow(bufferCapacity)
	for i, arg := range p.args {
		if i == 0 {
			buf.WriteString(fmt.Sprint(arg))
		} else {
			buf.WriteString(" ")
			buf.WriteString(fmt.Sprint(arg))
		}
	}
	return buf.String()

}

func (p *entry) WithContext(ctx context.Context) *entry {
	p.ctx = ctx
	return p
}

//func (p *entry) WithExtendField(key string, value any) *entry {
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
