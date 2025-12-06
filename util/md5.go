package util

import (
	"crypto/md5"
	"encoding/hex"
	xruntime "github.com/75912001/xlib/runtime"
	"github.com/pkg/errors"
	"io"
	"os"
)

// MD5 生成md5
func MD5(data []byte) string {
	sum := md5.Sum(data)
	return hex.EncodeToString(sum[:])
}

func MD5File(pathFile string) (md5sum string, err error) {
	f, err := os.Open(pathFile)
	if err != nil {
		return "", errors.WithMessagef(err, "open file %v %v", pathFile, xruntime.Location())
	}
	defer func() {
		_ = f.Close()
	}()
	md5hash := md5.New()
	_, err = io.Copy(md5hash, f)
	if err != nil {
		return "", errors.WithMessagef(err, "copy file %v %v", pathFile, xruntime.Location())
	}
	md5sum = hex.EncodeToString(md5hash.Sum(nil))
	return md5sum, nil
}
