package log

import (
	"bytes"
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
	if buf.Cap() < bufferCapacity {
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
	if p.ctx != nil {
		if traceIDVal := p.ctx.Value(TraceIDKey); traceIDVal != nil {
			if s, ok := traceIDVal.(string); ok {
				buf.WriteString(s)
			} else {
				_, _ = fmt.Fprintf(buf, "%v", traceIDVal)
			}
		} else {
			buf.WriteString("nil")
		}
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
	_, _ = fmt.Fprintf(buf, callerInfoFormat, p.line, p.file, p.funcName)
	buf.WriteByte(']')
	// 处理 fields, 转换为 json 格式
	buf.WriteString("[{")
	for idx, v := range p.extendFields {
		if idx%2 == 0 { // key
			buf.WriteByte('"')
			if str, ok := v.(string); ok {
				buf.WriteString(str)
			} else {
				_, _ = fmt.Fprintf(buf, "%v", v)
			}
			buf.WriteString(`":`)
		} else { // val
			buf.WriteByte('"')
			if str, ok := v.(string); ok {
				buf.WriteString(str)
			} else {
				_, _ = fmt.Fprintf(buf, "%v", v)
			}
			buf.WriteByte('"')
			if idx+1 < len(p.extendFields) {
				buf.WriteByte(',')
			}
		}
	}
	buf.WriteString("}]")
	// 日志消息(与主缓冲合并,避免 getMessage 二次缓冲与 String 分配)
	buf.WriteByte(' ')
	appendLogMessage(buf, p)

	p.outBytes = append(p.outBytes[:0], buf.Bytes()...)
}

// appendLogMessage 将用户消息写入 buf(与 formatLogData 共用同一缓冲)
func appendLogMessage(buf *bytes.Buffer, p *entry) {
	if p.format != "" {
		_, _ = fmt.Fprintf(buf, p.format, p.args...)
		return
	}
	for i, arg := range p.args {
		if i > 0 {
			_ = buf.WriteByte(' ')
		}
		_, _ = fmt.Fprintf(buf, "%v", arg)
	}
}
