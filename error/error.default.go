package error

import (
	"fmt"
	xmap "github.com/75912001/xlib/map"
	xpool "github.com/75912001/xlib/pool"
	"github.com/pkg/errors"
	"strconv"
)

type Error struct {
	code         uint32 // 码
	name         string // 名称
	desc         string // 描述 description
	extraMessage string // 附加信息
	extraError   error  // 附加错误
}

// NewError 创建错误码 [初始化程序的时候创建] 创建失败会 panic.
func NewError(code uint32) *Error {
	newObject := newError(code)
	_, ok := errMap.Find(newObject.code)
	if ok { // 重复
		panic(errors.WithStack(errors.Errorf("duplicate err, code:%v %#x", newObject.code, newObject.code)))
	}
	errMap.Add(code, struct{}{})
	return newObject
}

func (p *Error) Error() string {
	if Success.code == p.code {
		return ""
	}
	buf := xpool.Buffer.Get()
	defer func() {
		xpool.Buffer.Put(buf)
	}()
	buf.Grow(300)
	buf.WriteString("name:")
	buf.WriteString(p.name)
	buf.WriteString(" code:")
	buf.WriteString(strconv.FormatUint(uint64(p.code), 10))
	buf.WriteString(fmt.Sprintf(" %#x", p.code))
	buf.WriteString(" description:")
	buf.WriteString(p.desc)
	if p.extraMessage != "" {
		buf.WriteString(" extraMessage:")
		buf.WriteString(p.extraMessage)
	}
	if p.extraError != nil {
		buf.WriteString(" extraError:")
		buf.WriteString(p.extraError.Error())
	}
	outString := buf.String()
	return outString
}

func (p *Error) WithName(name string) *Error {
	p.name = name
	return p
}

func (p *Error) Name() string {
	return p.name
}

func (p *Error) WithDesc(desc string) *Error {
	p.desc = desc
	return p
}

func (p *Error) Desc() string {
	return p.desc
}

func (p *Error) WithExtraMessage(extraMessage string) *Error {
	p.extraMessage = extraMessage
	return p
}

func (p *Error) ExtraMessage() string {
	return p.extraMessage
}

func (p *Error) WithExtraError(extraError error) *Error {
	p.extraError = extraError
	return p
}

func (p *Error) ExtraError() error {
	return p.extraError
}

// 用来确保 错误码-唯一性
var errMap = xmap.NewMapMgr[uint32, struct{}]()

func newError(code uint32) *Error {
	return &Error{
		code: code,
	}
}
