package pkg

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/golang/glog"
	"k8s-dev/namespace-injector/pkg/common"
	"k8s-dev/namespace-injector/pkg/config"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	informers "k8s.io/client-go/informers/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
	"time"
)

const controllerAgentName = "namespace-injector"

type Controller struct {
	// kubeclientset is a standard kubernetes clientset
	kubeclientset           kubernetes.Interface
	configMapSynced         cache.InformerSynced
	configMapLister         v1.ConfigMapLister
	namespaceSynced         cache.InformerSynced
	namespaceLister         v1.NamespaceLister
	namespaceLabelSelector  labels.Selector


	workqueue        workqueue.RateLimitingInterface
	syncedConfigMaps map[string]string
}

func NewController(
	kubeclientset kubernetes.Interface,
	configMapInformer informers.ConfigMapInformer,
	namespaceInformer informers.NamespaceInformer,
) *Controller {
	controller := &Controller{
		kubeclientset:           kubeclientset,
		configMapSynced:         configMapInformer.Informer().HasSynced,
		configMapLister:         configMapInformer.Lister(),
		namespaceSynced:         namespaceInformer.Informer().HasSynced,
		namespaceLister:         namespaceInformer.Lister(),
		workqueue:               workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), controllerAgentName),
		syncedConfigMaps:        make(map[string]string),
	}
	if namespaceLabelSelector, err := labels.Parse(config.NamespaceLabelSelectors); err != nil {
		panic(err)
	} else {
		controller.namespaceLabelSelector = namespaceLabelSelector
	}
	// Set up an event handler for when ConfigMap resources change
	configMapInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: controller.enqueue,
		UpdateFunc: func(oldObj, newObj interface{}) {
			oldConfigMap := oldObj.(*corev1.ConfigMap)
			newConfigMap := newObj.(*corev1.ConfigMap)
			if oldConfigMap.ResourceVersion == newConfigMap.ResourceVersion {
				return
			}
			controller.enqueue(newObj)
		},
		DeleteFunc: controller.enqueueForDelete,
	})

	// Set up an event handler for when namespace resources change
	namespaceInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    controller.enqueue,
		UpdateFunc: func(oldObj, newObj interface{}) {},
		DeleteFunc: func(obj interface{}) {},
	})
	return controller
}

// Run will set up the event handlers for types we are interested in, as well
// as syncing informer caches and starting workers. It will block until stopCh
// is closed, at which point it will shutdown the workqueue and wait for
// workers to finish processing their current work items.
func (c *Controller) Run(threadiness int, stopCh <-chan struct{}) error {
	defer runtime.HandleCrash()
	defer c.workqueue.ShutDown()

	// Start the informer factories to begin populating the informer caches
	glog.Info("Starting Network control loop")

	// Wait for the caches to be synced before starting workers
	glog.Info("Waiting for informer caches to sync")
	if ok := cache.WaitForCacheSync(stopCh, c.configMapSynced, c.namespaceSynced); !ok {
		return fmt.Errorf("failed to wait for caches to sync")
	}

	glog.Info("Starting workers")
	// Launch two workers to process Network resources
	for i := 0; i < threadiness; i++ {
		go wait.Until(c.runWorker, time.Second, stopCh)
	}

	glog.Info("Started workers")
	<-stopCh
	glog.Info("Shutting down workers")

	return nil
}

// runWorker is a long-running function that will continually call the
// processNextWorkItem function in order to read and process a message on the
// workqueue.
func (c *Controller) runWorker() {
	for c.processNextWorkItem() {
	}
}

func (c *Controller) processNextWorkItem() bool {
	obj, shutdown := c.workqueue.Get()
	if shutdown {
		return false
	}
	// We wrap this block in a func so we can defer c.workqueue.Done.
	err := func(obj interface{}) error {
		// We call Done here so the workqueue knows we have finished
		// processing this item. We also must remember to call Forget if we
		// do not want this work item being re-queued. For example, we do
		// not call Forget if a transient error occurs, instead the item is
		// put back on the workqueue and attempted again after a back-off
		// period.
		defer c.workqueue.Done(obj)
		var key string
		var ok bool
		// We expect strings to come off the workqueue. These are of the
		// form namespace/name. We do this as the delayed nature of the
		// workqueue means the items in the informer cache may actually be
		// more up to date that when the item was initially put onto the
		// workqueue.
		if key, ok = obj.(string); !ok {
			// As the item in the workqueue is actually invalid, we call
			// Forget here else we'd go into a loop of attempting to
			// process a work item that is invalid.
			c.workqueue.Forget(obj)
			runtime.HandleError(fmt.Errorf("expected string in workqueue but got %#v", obj))
			return nil
		}
		// Run the syncHandler, passing it the namespace/name string of the
		// Network resource to be synced.
		if err := c.syncHandler(key); err != nil {
			return fmt.Errorf("error syncing '%s': %s", key, err.Error())
		}
		// Finally, if no error occurs we Forget this item so it does not
		// get queued again until another change happens.
		c.workqueue.Forget(obj)
		glog.Infof("Successfully synced '%s'", key)
		return nil
	}(obj)
	if err != nil {
		runtime.HandleError(err)
		return true
	}
	return true
}

func (c *Controller) syncHandler(key string) error {
	namespace, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		runtime.HandleError(fmt.Errorf("invalid resource key: %s", key))
		return nil
	}
	if namespace == "" {
		return c.syncNamespaceHandler(name)
	} else {
		return c.syncConfigMapHandler(namespace, name)
	}
}

func (c *Controller) syncNamespaceHandler(name string) error {
	namespaces, err := c.namespaceLister.List(c.namespaceLabelSelector)
	if err != nil {
		runtime.HandleError(fmt.Errorf("failed to list namespace by: %s", name))
		return err
	}
	namespace := common.GetNsFromList(namespaces, name)
	if namespace == nil {
		return nil
	}
	logrus.Infof("try to process namespace: %s", name)
	if files, err := common.ListAllConfigMapCachedFiles(); err != nil {
		logrus.Errorf("list all configMap files error: %s", err.Error())
	} else {
		common.CreateK8sResourcesInNS(files, name)
	}
	return nil
}

func (c *Controller) syncConfigMapHandler(namespace, name string) error {
	configMap, err := c.configMapLister.ConfigMaps(namespace).Get(name)
	if err != nil {
		if errors.IsNotFound(err) {
			c.removeConfigMap(namespace, name)
			return nil
		}
		runtime.HandleError(fmt.Errorf("failed to list configMap by: %s/%s", namespace, name))
		return err
	}
	if !common.IsValidConfigMap(configMap) {
		c.removeConfigMap(namespace, name)
	} else {
		c.addConfigMap(configMap)
	}
	return nil
}

func (c *Controller) removeConfigMap(namespace, name string){
	if _, ok := c.syncedConfigMaps[name]; !ok {
		return
	}
	logrus.Infof("deleting configMap: %s/%s", namespace, name)
	if files, err := common.ListConfigMapCachedFiles(name); err != nil {
		logrus.Errorf("list configMap %s files error: %s", name, err.Error())
	} else {
		namespaces, err := c.namespaceLister.List(c.namespaceLabelSelector)
		if err != nil {
			logrus.Error("list namespace error")
			return
		}
		for _, ns := range namespaces {
			common.DeleteK8sResourcesInNS(files, ns.Name)
		}
		common.RemoveFiles(files...)
	}
	delete(c.syncedConfigMaps, name)
}

func (c *Controller) addConfigMap(configMap *corev1.ConfigMap) {
	if resourceVersion, ok := c.syncedConfigMaps[configMap.Name]; ok && resourceVersion == configMap.ResourceVersion {
		return
	}
	logrus.Infof("try to process configmap: %s/%s", configMap.Namespace, configMap.Name)
	allFilePaths := make([]string, 0)
	for name, content := range configMap.Data {
		fileName := common.GetResourceFileName(configMap.Name, name)
		if filePath, err := common.WriteFile2Cache(fileName, content); err != nil {
			logrus.Errorf("write file %s error: %s", fileName, err.Error())
		} else {
			allFilePaths = append(allFilePaths, filePath)
		}
	}
	namespaces, err := c.namespaceLister.List(c.namespaceLabelSelector)
	if err != nil {
		logrus.Error("list namespace error")
		return
	}
	for _, namespace := range namespaces {
		common.CreateK8sResourcesInNS(allFilePaths, namespace.Name)
	}
	c.syncedConfigMaps[configMap.Name] = configMap.ResourceVersion
}

func (c *Controller) enqueue(obj interface{}) {
	var key string
	var err error
	if key, err = cache.MetaNamespaceKeyFunc(obj); err != nil {
		runtime.HandleError(err)
		return
	}
	c.workqueue.AddRateLimited(key)
}

func (c *Controller) enqueueForDelete(obj interface{}) {
	var key string
	var err error
	key, err = cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
	if err != nil {
		runtime.HandleError(err)
		return
	}
	c.workqueue.AddRateLimited(key)
}
