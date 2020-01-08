package etcd

import (
	"k8s.io/utils/exec"
	"os"
	"fmt"
	"strings"
	"errors"
	"sync"
)

const (
	EXECPATH = "/usr/bin/etcdctl"
	ENV_ETCDENDPOINTS = "ETCDENDPOINTS"
	ENV_ETCDCTL_CA_FILE = "ETCDCTL_CA_FILE"
	ENV_ETCDCTL_CERT_FILE = "ETCDCTL_CERT_FILE"
	ENV_ETCDCTL_KEY_FILE = "ETCDCTL_KEY_FILE"
)

type EtcdV2 struct {
	runner exec.Interface
	execPath string
	args []string
	lock *sync.RWMutex
}

func NewEtcdV2() *EtcdV2 {
	endpints := os.Getenv(ENV_ETCDENDPOINTS)
	cafile := os.Getenv(ENV_ETCDCTL_CA_FILE)
	certfile := os.Getenv(ENV_ETCDCTL_CERT_FILE)
	keyfile := os.Getenv(ENV_ETCDCTL_KEY_FILE)
	args := []string{
		fmt.Sprintf("--endpoints=%s", endpints),
	}
	if cafile != "" {
		args = append(args,
			fmt.Sprintf("--ca-file=%s", cafile),
			fmt.Sprintf("--cert-file=%s", certfile),
			fmt.Sprintf("--key-file=%s", keyfile))
	}
	return &EtcdV2{
		runner: exec.New(),
		execPath: EXECPATH,
		args:args,
		lock: new(sync.RWMutex),
	}
}

func (e EtcdV2)Get(key string) ([]byte, error) {
	e.lock.RLock()
	defer e.lock.RUnlock()
	args := make([]string, 0)
	args = append(args, e.args...)
	args = append(args, "get", key)
	cmd := e.runner.Command(e.execPath, args...)
	out, err := cmd.CombinedOutput()
	if strings.Contains(string(out), "Key not found"){
		return []byte{}, nil
	}
	return out, err
}

func (e EtcdV2) Set(key, value string) error {
	e.lock.Lock()
	defer e.lock.Unlock()
	args := make([]string, 0)
	args = append(args, e.args...)
	args = append(args, "set", key, value)
	cmd := e.runner.Command(e.execPath, args...)
	out, err := cmd.CombinedOutput()
	if err != nil{
		return errors.New(fmt.Sprintf("Error: %s\n, out: %s", err.Error(), string(out)))
	}
	return nil
}

