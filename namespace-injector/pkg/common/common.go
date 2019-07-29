package common

import (
	"fmt"
	"k8s-dev/namespace-injector/pkg/config"
	"k8s.io/api/core/v1"
	"os"
	"path/filepath"
	"strings"
)

func IsValidConfigMap(configMap *v1.ConfigMap) bool {
	if value, ok := configMap.Labels[config.ComfigMapLabel]; ok {
		return value == "open"
	}
	return false
}

func GetResourceFileName(configMapName, keyName string) string {
	return fmt.Sprintf("%s_%s", configMapName, keyName)
}

func SplitResourceFileName(fileName string) (string, string, error) {
	parts := strings.Split(fileName, "_")
	switch len(parts) {
	case 1:
		return "", parts[0], nil
	case 2:
		return parts[0], parts[1], nil
	}

	return "", "", fmt.Errorf("unexpected file name format: %s", fileName)
}

func ListConfigMapCachedFiles(configMapName string) ([]string, error) {
	return GetFilteredSubFiles(config.CacheDir, func(path string) bool {
		fileName := filepath.Base(path)
		if fileConfigMapName, _, err := SplitResourceFileName(fileName); err == nil && configMapName == fileConfigMapName {
			return true
		}
		return false
	})
}

func ListAllConfigMapCachedFiles() ([]string, error) {
	return GetFilteredSubFiles(config.CacheDir, func(path string) bool {
		fileName := filepath.Base(path)
		if fileConfigMapName, _, err := SplitResourceFileName(fileName); err == nil && fileConfigMapName != "" {
			return true
		}
		return false
	})
}

func GetNsFromList(namespaces []*v1.Namespace, name string) *v1.Namespace {
	for _, namespace := range namespaces {
		if namespace.Name == name {
			return namespace
		}
	}
	return nil
}

func MakeCacheDir() error {
	if exist, err := FileExists(config.CacheDir); err != nil {
		return err
	} else if exist {
		return nil
	}
	if err := os.MkdirAll(config.CacheDir, os.ModePerm); err != nil {
		return err
	}
	return nil
}