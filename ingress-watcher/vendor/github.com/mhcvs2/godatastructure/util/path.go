package util

import (
	"os"
	"path/filepath"
	"errors"
	"fmt"
)

// Exists 判断一个文件或文件夹是否存在
func Exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

//GetSubFiles 列出文件夹下的所有文件，不包括子文件夹下的文件
func GetSubFiles(root string) ([]string, error) {
	res := []string{}
	err := filepath.Walk(root, func(path string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}
		if path == root {
			if f.IsDir() {
				return nil
			} else {
				return errors.New(fmt.Sprintf("%s is not a directory", path))
			}
		}
		if f.IsDir(){
			return filepath.SkipDir
		}
		res = append(res, path)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}
