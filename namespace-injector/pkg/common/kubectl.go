package common

import (
	"fmt"
	"github.com/Sirupsen/logrus"
)

func CheckKubeCtl() {
	if err := RunShell("kubectl get node"); err != nil {
		panic(err)
	}
}

func CreateK8sResourceInNS(filePath, ns string) error {
	command := fmt.Sprintf("kubectl apply -f %s -n %s", filePath, ns)
	return RunShell(command)
}

func DeleteK8sResourceInNS(filePath, ns string) error {
	command := fmt.Sprintf("kubectl delete -f %s -n %s", filePath, ns)
	return RunShell(command)
}

func CreateK8sResourcesInNS(filePaths []string, ns string) {
	for _, filePath := range filePaths {
		if err := CreateK8sResourceInNS(filePath, ns); err != nil {
			logrus.Errorf("create k8s resource %s in ns %s error: %s", filePath, ns, err.Error())
		}
	}
}

func DeleteK8sResourcesInNS(filePaths []string, ns string) {
	for _, filePath := range filePaths {
		_ = DeleteK8sResourceInNS(filePath, ns)
	}
}


