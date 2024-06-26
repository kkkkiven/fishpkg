// Copyright (c) 2021. Homeland Interactive Technology Ltd. All rights reserved.

package util

import (
	"errors"
	"os"
	"path/filepath"
)

var ErrDirPath = errors.New("dir path error")

// FileExists 检查文件是否存在
func FileExists(path string) bool {
	_, err := os.Stat(path)
	if err != nil && os.IsNotExist(err) {
		return false
	}
	return true
}

// DirExists 检查目录是否存在
func DirExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return os.IsExist(err)
	} else {
		return info.IsDir()
	}
}

// CreateDir 递归创建目录
func CreateDir(name string) error {
	if DirExists(name) {
		return nil
	}

	// 分解上层目录
	pdir := filepath.Dir(name)
	if pdir == "" {
		return ErrDirPath
	}

	if !DirExists(pdir) {
		err := CreateDir(pdir)
		if err != nil {
			return err
		}
	}

	return os.Mkdir(name, 0750)
}
