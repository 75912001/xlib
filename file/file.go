// 文件

package file

import (
	xruntime "github.com/75912001/xlib/runtime"
	"github.com/pkg/errors"
	"os"
)

// WriteFile 写文件,可选择覆盖或者追加
//
//	pathFile: 文件路径
//	data: 写入的数据
//	opts: 覆盖写入, 追加写入, 默认覆盖写入
func WriteFile(pathFile string, data []byte, opts *options) error {
	var err error
	var file *os.File

	if opts.append {
		file, err = os.OpenFile(pathFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	} else if opts.overwrite {
		file, err = os.OpenFile(pathFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	}

	if err != nil {
		return errors.WithMessagef(err, "write file %v, %v", pathFile, xruntime.Location())
	}
	defer func() {
		_ = file.Close()
	}()

	_, err = file.Write(data)
	if err != nil {
		return errors.WithMessagef(err, "write file %v, %v", pathFile, xruntime.Location())
	}

	return nil
}

// PathFileExists 判断文件是否存在
//
//	pathFile: 文件路径
//	return: 存在-true, 不存在-false
func PathFileExists(pathFile string) bool {
	_, err := os.Stat(pathFile)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

// CreateDirectory 创建目录
//
//	path: 目录路径
func CreateDirectory(path string) error {
	_, err := os.Stat(path)
	if err == nil {
		return nil
	}
	if os.IsNotExist(err) {
		err = os.MkdirAll(path, 0755) //os.ModePerm
		if err != nil {
			return errors.WithMessagef(err, "create directory %v, %v", path, xruntime.Location())
		}
	}
	return nil
}
