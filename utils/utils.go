package utils

import (
	"crypto/md5"
	"encoding/hex"
	"os"
	"path/filepath"
)

// is folder or file exist
func Exists(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

// is folder or not
func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

// is file or not
func IsFile(path string) bool {
	return !IsDir(path)
}

func GetDirSize(rootPath string) (uint64, error) {
	dirSize := uint64(0)

	readSize := func(path string, file os.FileInfo, err error) error {
		if err == nil && file != nil && !file.IsDir() {
			dirSize += uint64(file.Size())
		}

		return nil
	}

	err := filepath.Walk(rootPath, readSize)
	return dirSize, err
}

func GetStringHash(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}
