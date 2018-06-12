package main

import (
	"k8s.io/utils/exec"
	"fmt"
	"os"
	"encoding/json"
	"strings"
)

func t1() {
	os.Setenv("http_proxy", "")
	os.Setenv("https_proxy", "")
	runner := exec.New()
	execPath := "/usr/bin/etcdctl"
	args := []string{}
	args = append(args, "--endpoint=https://109.105.1.253:2379",
		                      "--ca-file=/etc/etcd/ssl/etcd-ca.pem",
		                      "--cert-file=/etc/etcd/ssl/etcd.pem",
		                      "--key-file=/etc/etcd/ssl/etcd-key.pem")

	args = append(args, "get")
	args = append(args, "/confd-dns/gpu-infra/dns/names")
	cmd := runner.Command(execPath, args...)
	out, err := cmd.CombinedOutput()
	notFound := false
	if err != nil {
		if strings.Contains(string(out), "Key not found"){
			notFound = true
		} else {
			fmt.Println(string(out))
			panic(err)
		}
	}
	names := []string{}
	if ! notFound{
		fmt.Println("unmarshal")
		err = json.Unmarshal(out, &names)
		if err != nil {
			panic(err)
		}
	}
	if len(names) == 0 {
		fmt.Println("kong...")
	} else {
		for _, name := range names {
			fmt.Println(name)
		}
	}
}

func main() {
	t1()
}