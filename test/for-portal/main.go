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
	batchv1 "k8s.io/api/batch/v1"
	"regexp"
	"strings"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"
	"strconv"
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
	quotaInformer := sharedInformerFactory.Core().V1().ResourceQuotas().Informer()
	go podInformer.Run(stopCh)
	go serviceInformer.Run(stopCh)
	go jobInformer.Run(stopCh)
	go quotaInformer.Run(stopCh)

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
	quotaLister := sharedInformerFactory.Core().V1().ResourceQuotas().Lister()
	fmt.Println(podLister)
	fmt.Println(serviceLister)
	fmt.Println(jobLister)
	fmt.Println(quotaLister)

	////获取一个开发环境的ssh端口
	//sshPortRequirement, _ := labels.NewRequirement("release", selection.Equals, []string{"hongchao.ma-dev1"})
	//sshPortSelector := labels.NewSelector().Add(*sshPortRequirement)
	//services, _ := serviceLister.Services("test-lab").List(sshPortSelector)
	//fmt.Println(services[0].Spec.Ports[0].NodePort)
	//
	////获取lab内所有开发容器
	//devSituationRequirement, _ := labels.NewRequirement("app", selection.Equals, []string{"dev-situation"})
	//devSituationSelector := labels.NewSelector().Add(*devSituationRequirement)
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
	//devSituationRequirement, _ := labels.NewRequirement("app", selection.Equals, []string{"dev-situation"})
	//devSituationRequirement2, _ := labels.NewRequirement("user", selection.Equals, []string{"hongchao-ma"})
	//devSituationSelector := labels.NewSelector().Add(*devSituationRequirement, *devSituationRequirement2)
	//devPods, err := podLister.Pods("test-lab").List(devSituationSelector)
	//var containerStatus v1.ContainerStatus
	//for i, dp := range devPods {
	//	containerStatus = dp.Status.ContainerStatuses[0]
	//	fmt.Printf("第%d个dev situation---------------\n", i)
	// if dp.Status.StartTime != nil {
	//	   fmt.Printf("Start time: %s\n", dp.Status.StartTime.String())
	//}
	//	fmt.Printf("image name is %s\n", containerStatus.Image)
	//	fmt.Printf("state is %s\n", GetContainerStatus(containerStatus.State))
	//}

	//获取指定用户运行任务的job
	//runJobRequirement, _ := labels.NewRequirement("app", selection.In, []string{"run-job"})
	//runJobRequirement2, _ := labels.NewRequirement("user", selection.In, []string{"admin"})
	//runJobSelector := labels.NewSelector().Add(*runJobRequirement, *runJobRequirement2)
	//runJobs, _ := jobLister.Jobs("intelligence-data-lab").List(runJobSelector)
	//for i, rj := range runJobs {
	//	fmt.Printf("第%d个job---------------\n", i)
	//	if rj.Status.StartTime != nil {
	//		fmt.Printf("Start time: %s\n", rj.Status.StartTime.String())
	//	}
	//
	//	fmt.Printf("image name is %s\n", rj.Spec.Template.Spec.Containers[0].Image)
	//	fmt.Printf("state is %s\n", GetJobStatus(rj.Status))
	//}

	//获取指定用户运行任务的job
	//runJobRequirement, _ := labels.NewRequirement("app", selection.In, []string{"commit-job"})
	////runJobRequirement2, _ := labels.NewRequirement("user", selection.In, []string{"admin"})
	//runJobRequirement2, _ := labels.NewRequirement("user", selection.In, []string{"hongchao-ma"})
	//runJobSelector := labels.NewSelector().Add(*runJobRequirement, *runJobRequirement2)
	//runJobs, _ := jobLister.Jobs("intelligence-data-lab").List(runJobSelector)
	//for i, rj := range runJobs {
	//	fmt.Printf("第%d个job---------------\n", i)
	//	if rj.Status.StartTime != nil {
	//		fmt.Printf("Start time: %s\n", rj.Status.StartTime.String())
	//	}
	//
	//	fmt.Printf("image name is %s\n", rj.Spec.Template.Spec.Containers[0].Image)
	//	fmt.Printf("state is %s\n", GetJobStatus(rj.Status))
	//}

	//获取resource quota
	//获取所有quota
	//quotas, _ := quotaLister.List(labels.Everything())
	//for _, quota := range quotas {
	//	fmt.Printf("lab: %s\n",quota.ObjectMeta.Namespace)
	//	for k, v := range quota.Status.Used {
	//		total, _ := quota.Status.Hard[k]
	//		fmt.Printf("resource name: %s, total %s,  use %s\n", GetResourceName(string(k)), total.String(), v.String())
	//	}
	//}
	//获取指定lab 的 quota
	//quotas, _ := quotaLister.ResourceQuotas("intelligence-data-lab").List(labels.Everything())
	//if len(quotas) > 0 {
	//	quota := quotas[0]
	//	for k, v := range quota.Status.Used {
	//		total, _ := quota.Status.Hard[k]
	//		fmt.Printf("resource name: %s, total %s,  use %s\n", GetResourceName(string(k)), total.String(), v.String())
	//	}
	//}

	//修改lab的quota
	//quotas, _ := quotaLister.ResourceQuotas("intelligence-data-lab").List(labels.Everything())
	//if len(quotas) > 0 {
	//	quota := quotas[0]
	//	for k, v := range quota.Status.Used {
	//		total, _ := quota.Status.Hard[k]
	//		fmt.Printf("resource name: %s, total %s,  use %s\n", GetResourceName(string(k)), total.String(), v.String())
	//	}
	//
	//	fmt.Println("---------------------------------")
	//	resourceName := "k20m"
	//	newQuantity := 5
	//	fmt.Println(quota.Spec.Hard)
	//	exactResourceName, _ := GetExactKeyName(resourceName, quota.Spec.Hard)
	//	quantity := resource.Quantity{}
	//	quantity.Set(int64(newQuantity))
	//	quota.Spec.Hard[exactResourceName] = quantity
	//
	//	_, err = clientset.CoreV1().ResourceQuotas("intelligence-data-lab").Update(quota)
	//	if err != nil {
	//		panic(err)
	//	}
	//}

	////获取每个lab的gpu资源使用情况(正在使用的)
	//devSituationRequirement, _ := labels.NewRequirement("type", selection.Equals, []string{"gpu-infra"})
	//devSituationSelector := labels.NewSelector().Add(*devSituationRequirement)
	//devPods, err := podLister.Pods("machine-learning-lab").List(devSituationSelector)
	//var containerStatus v1.ContainerStatus
	//for _, dp := range devPods {
	//	containerStatus = dp.Status.ContainerStatuses[0]
	//	if containerStatus.Ready {
	//		user, _ := dp.ObjectMeta.Labels["user"]
	//		singleId := strings.Replace(user, "-", ".", -1)
	//		fmt.Println("user: ", singleId)
	//		for k, v := range dp.Spec.Containers[0].Resources.Limits {
	//			number, _ := strconv.Atoi(v.String())
	//			fmt.Printf("podName: %s, gpu type: %s, number is %d\n", dp.Name, GetResourceName(k.String()), number)
	//		}
	//	}
	//}

	//获取每个lab的gpu资源使用情况
	devSituationRequirement, _ := labels.NewRequirement("type", selection.Equals, []string{"gpu-infra"})
	devSituationSelector := labels.NewSelector().Add(*devSituationRequirement)
	devPods, err := podLister.Pods("machine-learning-lab").List(devSituationSelector)
	for _, dp := range devPods {
		user, _ := dp.ObjectMeta.Labels["user"]
		singleId := strings.Replace(user, "-", ".", -1)
		fmt.Println("user: ", singleId)
		for k, v := range dp.Spec.Containers[0].Resources.Limits {
			number, _ := strconv.Atoi(v.String())
			fmt.Printf("podName: %s, gpu type: %s, number is %d\n", dp.Name, GetResourceName(k.String()), number)
		}
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

func GetResourceName(origin string) string {
	return regexp.MustCompile("cpu|memory|gpu.*").FindString(origin)
}

func GetExactKeyName(simpleName string, list v1.ResourceList) (v1.ResourceName, error) {
	for k := range list {
		if strings.Contains(k.String(), simpleName) {
			return k, nil
		}
	}
	return "", fmt.Errorf("resource %s not found", simpleName)
}


