package utils

import (
	"bytes"
	"errors"
	"io/ioutil"
	"os"
	"runtime"
)

// FileIsExist 检查文件是否存在
func FileIsExist(path string) bool {
	_, err := os.Stat(path)
	if err != nil && os.IsNotExist(err) {
		return false
	}
	return true
}

// GetFileContents 读取文件内容
func GetFileContents(path string) string {
	fi, err := os.Open(path)
	if err != nil {
		return ""
	}
	defer fi.Close()
	fd, err := ioutil.ReadAll(fi)
	return string(bytes.Trim(fd, "\xef\xbb\xbf"))
}

// DirIsExist 检查目录是否存在
func DirIsExist(path string) bool {
	fileinfo, err := os.Stat(path)
	if err != nil {
		return os.IsExist(err)
	} else {
		return fileinfo.IsDir()
	}
}

// 打开文件
func OpenFile(path string, fileName string, openType int) (*os.File, error) {

	fileFullPath := path + "/" + fileName
	if runtime.GOOS == "windows" {
		fileFullPath = path + "\\" + fileName
	}
	isFileExist := FileIsExist(fileFullPath)

	switch openType {
	case FILE_OPEN_TYPE_READ:
		if isFileExist {
			f, err := os.OpenFile(fileFullPath, os.O_RDONLY, 0644)
			return f, err
		} else {
			return nil, errors.New("file not exist")
		}
	case FILE_OPEN_TYPE_WRITE:
		if isFileExist {
			return nil, errors.New("file already existed")
		} else {
			f, err := os.OpenFile(fileFullPath, os.O_CREATE|os.O_WRONLY, 0644)
			return f, err
		}
	case FILE_OPEN_TYPE_APPEND:
		if !isFileExist {
			f, err := os.OpenFile(fileFullPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
			return f, err
		} else {
			f, err := os.OpenFile(fileFullPath, os.O_APPEND|os.O_WRONLY, 0644)
			return f, err
		}
	}
	return nil, nil
}
