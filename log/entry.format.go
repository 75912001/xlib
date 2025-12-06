package log

import (
	"fmt"
	xpool "github.com/75912001/xlib/pool"
	"strconv"
)

// 格式化日志数据
// 格式为  [时间][日志级别][TID:xxx][UID:xxx][堆栈信息][扩展信息,为 json 格式 {ExtendFields-key:ExtendFields:val,...}][日志消息]

func formatLogData(p *entry) {
	// 使用sync.Pool来复用buffer以减少内存分配
	buf := xpool.Buffer.Get()
	defer func() {
		xpool.Buffer.Put(buf)
	}()
	capacity := buf.Cap()
	if capacity < bufferCapacity {
		buf.Grow(bufferCapacity)
	}

	// 时间
	buf.WriteByte('[')
	buf.WriteString(p.time.Format(logTimeFormat))
	buf.WriteByte(']')
	// 日志级别
	buf.WriteByte('[')
	buf.WriteString(levelDesc[p.level])
	buf.WriteByte(']')
	// TraceID
	buf.WriteByte('[')
	buf.WriteString(TraceIDKey)
	buf.WriteByte(':')
	if p.ctx != nil { // 处理 ctx 中的 traceID
		if traceIdVal := p.ctx.Value(TraceIDKey); traceIdVal != nil {
			buf.WriteString(traceIdVal.(string))
		} else {
			buf.WriteString("nil")
		}
	} else { // 没有 ctx , 则 traceID 为 ""
	}
	buf.WriteByte(']')
	// UID 优先从 ctx 查找,其次查找 field 当 UID 不存在时也需要占位 值为0
	var uid uint64
	if p.ctx != nil {
		if uidVal := p.ctx.Value(UserIDKey); uidVal != nil {
			uid, _ = uidVal.(uint64)
		}
	}
	if uid == 0 { //没有找到UID,从fields中查找,找到第一个
		for idx := 0; idx+1 < len(p.extendFields); idx += 2 {
			if key, ok := p.extendFields[idx].(string); ok && key == UserIDKey { //找到
				uid, _ = p.extendFields[idx+1].(uint64)
				break
			}
		}
	}
	buf.WriteByte('[')
	buf.WriteString(UserIDKey)
	buf.WriteByte(':')
	buf.WriteString(strconv.FormatUint(uid, 10))
	buf.WriteByte(']')
	// 堆栈
	buf.WriteByte('[')
	buf.WriteString(p.getCallerInfo())
	buf.WriteByte(']')
	// 处理 fields, 转换为 json 格式
	buf.WriteString(fmt.Sprint("[{"))
	for idx, v := range p.extendFields {
		if idx%2 == 0 { //key
			buf.WriteByte('"')
			if str, ok := v.(string); ok {
				buf.WriteString(str)
			} else {
				buf.WriteString(fmt.Sprint(v))
			}
			buf.WriteString(`":`)
		} else { //val
			buf.WriteByte('"')
			if str, ok := v.(string); ok {
				buf.WriteString(str)
			} else {
				buf.WriteString(fmt.Sprint(v))
			}
			buf.WriteByte('"')
			if idx+1 < len(p.extendFields) {
				buf.WriteByte(',')
			}
		}
	}
	buf.WriteString(fmt.Sprint("}]"))
	// 日志消息
	buf.WriteByte(' ')
	buf.WriteString(p.getMessage())
	p.outString = buf.String()
}
