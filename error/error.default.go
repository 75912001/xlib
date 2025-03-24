package error

import (
	"fmt"
	"github.com/pkg/errors"
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
	e := checkDuplication(newObject)
	if e != nil {
		panic(e)
	}
	errMap[code] = struct{}{}
	return newObject
}

func (p *Error) Error() string {
	if Success.code == p.code {
		return ""
	}
	return fmt.Sprintf("name:%v code:%v %#x description:%v extraMessage:%v extraError:%v",
		p.name, p.code, p.code, p.desc, p.extraMessage, p.extraError)
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
var errMap = make(map[uint32]struct{})

// 检查重复情况
func checkDuplication(err *Error) error {
	if _, ok := errMap[err.code]; ok { // 重复
		return errors.WithStack(errors.Errorf("duplicate err, code:%v %#x", err.code, err.code))
	}
	return nil
}

func newError(code uint32) *Error {
	return &Error{
		code: code,
	}
}
