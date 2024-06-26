package utils

import (
	"errors"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"fmt"
)

const (
	FILE_OPEN_TYPE_READ   int = 0
	FILE_OPEN_TYPE_WRITE  int = 1
	FILE_OPEN_TYPE_APPEND int = 2
)

// ReadTextFile 读取文件
func ReadTextFile(path string) string {
	fmt.Println("path:",path)
	fp, err := os.Open(path)
	if nil != err {
		return ""
	}
	defer fp.Close()
	d, err := ioutil.ReadAll(fp)
	if nil != err {
		return ""
	}
	return string(d)
}

/*
 * 读取配置文件
 */
func ReadConfFile(file string) string {
	return ReadTextFile(GetConfPath(file))
}

// 获取配置文件路径
func GetConfPath(fname string) string {
	// 绝对路径
	if len(fname) > 2 && ('/' == fname[0] || ':' == fname[1]) {
		return fname
	}
	return GetFilePathBaseExe("./" + fname)
}

// 以exe为根节点获取文件路径
func GetFilePathBaseExe(fname string) string {
	// 绝对路径
	if len(fname) > 2 && ('/' == fname[0] || ':' == fname[1]) {
		return fname
	}
	// 相对路径
	file, _ := exec.LookPath(os.Args[0])
	return filepath.Join(filepath.Dir(file), fname)
}

// 创建目录
func CreateDir(dirName string) error {
	// 判断目录是否存在
	if _, err := os.Stat(dirName); false == os.IsNotExist(err) {
		return nil
	}

	// 分解上层目录
	pdir := filepath.Dir(dirName)
	if "" == pdir {
		return errors.New("error dir format")
	}

	// 判断目录是否存在
	if _, err := os.Stat(pdir); os.IsNotExist(err) {
		err = CreateDir(pdir)
		if nil != err {
			return err
		}
	}

	return os.Mkdir(dirName, 0755)
}
