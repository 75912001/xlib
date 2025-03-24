package log

import (
	"fmt"
	xruntime "github.com/75912001/xlib/runtime"
	"github.com/pkg/errors"
	"os"
)

// 生成 normal log Writer
func newNormalFileWriter(filePath string, namePrefix string, logDuration int) (*os.File, error) {
	return newFileWriter(filePath, namePrefix, logDuration, normalLogFileBaseName)
}

// 生成 error log Writer
func newErrorFileWriter(filePath string, namePrefix string, logDuration int) (*os.File, error) {
	return newFileWriter(filePath, namePrefix, logDuration, errorLogFileBaseName)
}

// 生成 log Writer
func newFileWriter(filePath string, namePrefix string, logDuration int, fileBaseName string) (*os.File, error) {
	fileName := fmt.Sprintf(fileFormat, filePath, namePrefix, logDuration, fileBaseName)
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_APPEND|os.O_RDWR, os.FileMode(0644))
	if err != nil {
		return nil, errors.WithMessage(err, xruntime.Location())
	}
	return file, nil
}
