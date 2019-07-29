package common

import (
	"errors"
	"fmt"
	"github.com/Sirupsen/logrus"
	"io/ioutil"
	"k8s-dev/namespace-injector/pkg/config"
	"os"
	"path"
	"path/filepath"
)

func WriteFile(name, content string) error {
	data := []byte(content)
	if err := ioutil.WriteFile(name, data, 0644); err != nil {
		return err
	}
	logrus.Infof("write file %s success", name)
	return nil
}

func WriteFile2Cache(name, content string) (string, error) {
	filePath := GetCachedFilePath(name)
	return filePath, WriteFile(filePath, content)
}

func GetCachedFilePath(name string) string {
	return path.Join(config.CacheDir, name)
}

func HomeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE")
}

func FileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func GetFilteredSubFiles(root string, filterFunc func(path string) bool) ([]string, error) {
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
		if filterFunc(path) {
			res = append(res, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func RemoveFile(file string) {
	if exist, err := FileExists(file); err == nil && exist {
		if err = os.Remove(file); err != nil {
			logrus.Errorf("remove file %s error: %s", file, err.Error())
		}
	}
}

func RemoveFiles(files ...string) {
	for _, file := range files {
		RemoveFile(file)
	}
}
