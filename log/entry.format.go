package log

import (
	"bytes"
	"fmt"
	"strconv"
)

// 格式化日志数据
// 格式为  [时间][日志级别][TID:xxx][UID:xxx][堆栈信息][扩展信息,为 json 格式 {ExtendFields-key:ExtendFields:val,...}][日志消息]
func formatLogData(p *entry) {
	var buf bytes.Buffer
	buf.Grow(bufferCapacity)
	// 时间
	buf.WriteString(fmt.Sprint("[", p.time.Format(logTimeFormat), "]"))
	// 日志级别
	buf.WriteString(fmt.Sprint("[", levelDesc[p.level], "]"))
	// TraceID
	if p.ctx != nil { // 处理 ctx 中的 traceID
		traceIdVal := p.ctx.Value(TraceIDKey)
		if traceIdVal != nil {
			buf.WriteString(fmt.Sprint("[", TraceIDKey, ":", traceIdVal.(string), "]"))
		} else {
			buf.WriteString(fmt.Sprint("[", TraceIDKey, ":traceIdVal is nil]"))
		}
	} else { // 没有 ctx , 则 traceID 为0
		buf.WriteString(traceIDKeyString0)
	}
	// UID 优先从 ctx 查找,其次查找 field 当 UID 不存在时也需要占位 值为0
	var uid uint64
	if p.ctx != nil {
		uidVal := p.ctx.Value(UserIDKey)
		if uidVal != nil {
			uid, _ = uidVal.(uint64)
		}
	}
	if 0 == uid { //没有找到UID,从fields中查找,找到第一个
		for idx, v := range p.extendFields {
			str, ok := v.(string)
			if ok && str == UserIDKey { //找到
				uid, _ = p.extendFields[idx+1].(uint64)
				break
			}
		}
	}
	buf.WriteString(fmt.Sprint("[", UserIDKey, ":", strconv.FormatUint(uid, 10), "]"))
	// 堆栈
	buf.WriteString(fmt.Sprint("[", p.callerInfo, "]"))
	// 处理 fields, 转换为 json 格式
	buf.WriteString(fmt.Sprint("[{"))
	for idx, v := range p.extendFields {
		if idx%2 == 0 { //key
			//buf.WriteString("{")
			str, ok := v.(string)
			if ok {
				buf.WriteString(`"` + str + `"`)
			} else {
				buf.WriteString(fmt.Sprintf("\"%v\"", v))
			}
			buf.WriteString(":")
		} else { //val
			str, ok := v.(string)
			if ok {
				buf.WriteString(`"` + str + `"`)
			} else {
				buf.WriteString(fmt.Sprintf("\"%v\"", v))
			}
			if idx+1 < len(p.extendFields) {
				buf.WriteString(",")
			}
		}
	}
	buf.WriteString(fmt.Sprint("}]"))
	// 日志消息
	buf.WriteString(p.message)
	p.outString = buf.String()
}
