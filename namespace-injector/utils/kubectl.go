package utils

import "fmt"

func init() {
	if err := RunShell("kubectl get node"); err != nil {
		panic(err)
	}
}

func CreateK8sResourceInNS(filePath, ns string) error {
	command := fmt.Sprintf("kubectl apply -f %s -n %s", filePath, ns)
	return RunShell(command)
}


