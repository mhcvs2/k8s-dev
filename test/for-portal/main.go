package main

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"fmt"
	"path/filepath"
	"os"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/informers"
	"time"
	"k8s.io/client-go/tools/cache"
	"k8s.io/apimachinery/pkg/util/runtime"
	"os/signal"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"
	batchv1 "k8s.io/api/batch/v1"
)

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE")
}

func main() {
	signalCh := make(chan os.Signal, 1)
	go signal.Notify(signalCh, os.Interrupt, os.Kill)
	config, err := rest.InClusterConfig()
	if err != nil {
		home := homeDir()
		kubeconfig := filepath.Join(home, ".kube", "config")
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			panic(err.Error())
		}
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	sharedInformerFactory := informers.NewSharedInformerFactory(clientset, time.Minute*10)
	stopCh := make(chan struct{})
	podInformer := sharedInformerFactory.Core().V1().Pods().Informer()
	serviceInformer := sharedInformerFactory.Core().V1().Services().Informer()
	jobInformer := sharedInformerFactory.Batch().V1().Jobs().Informer()
	go podInformer.Run(stopCh)
	go serviceInformer.Run(stopCh)
	go jobInformer.Run(stopCh)

	if !cache.WaitForCacheSync(stopCh,
		podInformer.HasSynced,
		serviceInformer.HasSynced,
		jobInformer.HasSynced,
		) {
		runtime.HandleError(fmt.Errorf("timed out waiting for caches to sync"))
	} else {
		fmt.Println("success sync")
	}

	podLister := sharedInformerFactory.Core().V1().Pods().Lister()
	serviceLister := sharedInformerFactory.Core().V1().Services().Lister()
	jobLister := sharedInformerFactory.Batch().V1().Jobs().Lister()
	fmt.Println(podLister)
	fmt.Println(serviceLister)
	fmt.Println(jobLister)

	//pods1, err := podLister.List(labels.NewSelector())
	//if err != nil {
	//	panic(err.Error())
	//}
	//fmt.Println(len(pods1))
	//for _, p := range pods1 {
	//	fmt.Printf("%s ", p.Name)
	//}
	//fmt.Println()

	//pods2, err := podLister.Pods("kube-system").List(labels.Everything())
	//if err != nil {
	//	panic(err.Error())
	//}
	//fmt.Println(len(pods2))
	//for _, p := range pods2 {
	//	fmt.Printf("%s ", p.Name)
	//}
	//fmt.Println()
	//p, err := podLister.Pods("kube-system").Get("haproxy-k8s-ceph5")
	//if err != nil {
	//	panic(err.Error())
	//}
	//fmt.Println(p.Name)

	//获取一个开发环境的ssh端口
	//sshPortSelector := labels.NewSelector()
	//sshPortRequirement, _ := labels.NewRequirement("release", selection.Equals, []string{"hongchao.ma-dev1"})
	//sshPortSelector.Add(*sshPortRequirement)
	//services, _ := serviceLister.Services("test-lab").List(sshPortSelector)
	//fmt.Println(services[0].Spec.Ports[0].NodePort)

	//获取lab内所有开发容器
	//devSituationSelector := labels.NewSelector()
	//devSituationRequirement, _ := labels.NewRequirement("app", selection.Equals, []string{"dev-situation"})
	//devSituationSelector.Add(*devSituationRequirement)
	//devPods, err := podLister.Pods("test-lab").List(devSituationSelector)
	//var containerStatus v1.ContainerStatus
	//for i, dp := range devPods {
	//	containerStatus = dp.Status.ContainerStatuses[0]
	//	fmt.Printf("第%d个dev situation---------------\n", i)
	//	fmt.Printf("Start time: %s\n", dp.Status.StartTime.String())
	//	fmt.Printf("image name is %s\n", containerStatus.Image)
	//	fmt.Printf("state is %s\n", GetContainerStatus(containerStatus.State))
	//}

	//获取指定用户的开发容器
	//devSituationSelector := labels.NewSelector()
	//devSituationRequirement, _ := labels.NewRequirement("app", selection.Equals, []string{"dev-situation"})
	//devSituationRequirement2, _ := labels.NewRequirement("user", selection.Equals, []string{"hongchao-ma"})
	//devSituationSelector.Add(*devSituationRequirement, *devSituationRequirement2)
	//devPods, err := podLister.Pods("test-lab").List(devSituationSelector)
	//var containerStatus v1.ContainerStatus
	//for i, dp := range devPods {
	//	containerStatus = dp.Status.ContainerStatuses[0]
	//	fmt.Printf("第%d个dev situation---------------\n", i)
	//  if dp.Status.StartTime != nil {
	//	   fmt.Printf("Start time: %s\n", dp.Status.StartTime.String())
	// }
	//	fmt.Printf("image name is %s\n", containerStatus.Image)
	//	fmt.Printf("state is %s\n", GetContainerStatus(containerStatus.State))
	//}

	//获取指定用户运行任务的job
	runJobRequirement, _ := labels.NewRequirement("app", selection.In, []string{"run-job"})
	runJobRequirement2, _ := labels.NewRequirement("user", selection.In, []string{"admin"})
	runJobSelector := labels.NewSelector().Add(*runJobRequirement, *runJobRequirement2)
	runJobs, _ := jobLister.Jobs("intelligence-data-lab").List(runJobSelector)
	for i, rj := range runJobs {
		fmt.Printf("第%d个job---------------\n", i)
		if rj.Status.StartTime != nil {
			fmt.Printf("Start time: %s\n", rj.Status.StartTime.String())
		}

		fmt.Printf("image name is %s\n", rj.Spec.Template.Spec.Containers[0].Image)
		fmt.Printf("state is %s\n", GetJobStatus(rj.Status))
	}

	select {
	case <-stopCh:
		os.Exit(0)
	case <-signalCh:
		fmt.Println("exit by signal")
		os.Exit(0)
	}
}

func GetContainerStatus(state v1.ContainerState) string {
	if state.Running != nil {
		return "running"
	}
	if state.Waiting != nil {
		return "waiting"
	}
	if state.Terminated != nil {
		return "terminated"
	}
	return "unknow"
}

func GetJobStatus(state batchv1.JobStatus) string {
	if state.Active > 0 {
		return "active"
	}
	if state.Succeeded > 0 {
		return "succeeded"
	}
	if state.Failed > 0 {
		return "failed"
	}
	return "wait"
}
