package watcher

import (
	"encoding/json"
	"fmt"
	"github.com/mhcvs2/godatastructure/set"
	myetcd "k8s-dev/ingress-watcher/etcd"
	"k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"os"
	"time"
	"regexp"
	log "github.com/sirupsen/logrus"
)

const (
	ENV_WATCH_KEY = "WATCH_KEY"
	ENV_BASE_DOMAIN = "BASE_DOMAIN"
)

type IngressWatcher struct {
	etcd       *myetcd.EtcdV2
	watchKey   string
	informer   cache.SharedIndexInformer
	data       *set.HashSet
	baseDomain string
}

func NewIngressWatcher(clientset *kubernetes.Clientset) *IngressWatcher {
	sharedInformerFactory := informers.NewSharedInformerFactory(clientset, time.Minute*10)
	ingInformer := sharedInformerFactory.Extensions().V1beta1().Ingresses().Informer()
	watchKey := os.Getenv(ENV_WATCH_KEY)
	etcdcli := myetcd.NewEtcdV2()
	return &IngressWatcher{
		etcd:       etcdcli,
		watchKey:   watchKey,
		informer:   ingInformer,
		data:       set.NewHashSet(),
		baseDomain: os.Getenv(ENV_BASE_DOMAIN),
	}
}

func (ing *IngressWatcher) getData() error {
	out, err := ing.etcd.Get(ing.watchKey)
	if err != nil {
		return fmt.Errorf("Update data error: %s\n", err.Error())
	}
	tmpData := []string{}
	if err = json.Unmarshal(out, &tmpData); err != nil {
		return fmt.Errorf("data unmarshal error: %s\n", err.Error())
	}
	for _, d := range tmpData {
		ing.data.Add(d)
	}
	log.Infof("data: %v", ing.data.Elements())
	return nil
}

func (ing *IngressWatcher) addData(names ...string) bool {
	change := false
	for _, name := range names {
		if name == "" {
			continue
		}
		if !ing.data.Contains(name) {
			ing.data.Add(name)
			change = true
			log.Infof("Add domain %s\n", name)
		} else {
			log.Debugf("name %s exist, ignore it\n", name)
		}
	}
	return change
}

func (ing *IngressWatcher) deleteData(names ...string) bool {
	change := false
	for _, name := range names {
		if name == "" {
			continue
		}
		if ing.data.Contains(name) {
			ing.data.Remove(name)
			change = true
			log.Infof("Delete domain %s\n", name)
		} else {
			log.Debugf("name %s does exist, ignore it\n", name)
		}
	}
	return change
}

func (ing *IngressWatcher) getDomainName(host string) string {
	reg := regexp.MustCompile(fmt.Sprintf("%s$", ing.baseDomain))
	endWith := reg.MatchString(host)
	reg = regexp.MustCompile(fmt.Sprintf("^%s", ing.baseDomain))
	startWith := reg.MatchString(host)
	if startWith && endWith {
		log.Infof("Ignore base host %s\n", host)
	} else if endWith {
		reg = regexp.MustCompile(fmt.Sprintf(".%s", ing.baseDomain))
		return reg.ReplaceAllString(host, "")
	} else {
		log.Infof("Ignore host %s\n", host)
	}
	return ""
}

func (ing *IngressWatcher) Run(stopCh chan struct{}) {
	log.Info("Start Ingress Watcher...")
	if err := ing.getData(); err != nil {
		log.Errorln(err.Error())
	}
	go ing.informer.Run(stopCh)
	if !cache.WaitForCacheSync(stopCh,
		ing.informer.HasSynced,
	) {
		runtime.HandleError(fmt.Errorf("Timed out waiting for caches to sync"))
	} else {
		log.Info("success sync")
	}
	ingEventHandler := cache.ResourceEventHandlerFuncs{
		AddFunc:    ing.add,
		DeleteFunc: ing.del,
		UpdateFunc: ing.update,
	}
	ing.informer.AddEventHandler(ingEventHandler)
}

func (i *IngressWatcher) updateEtcd(eles []interface{}) {
	if out, err := json.Marshal(eles); err != nil {
		log.Errorf("Marshal error: %s\n", err.Error())
		log.Errorf("Broken data is: %v\n", eles)
	} else if err = i.etcd.Set(i.watchKey, string(out)); err != nil {
		log.Error(err.Error())
	} else {
		log.Infof("update: %s", string(out))
	}
}

//-----------------------------------------------

func (i *IngressWatcher) add(obj interface{}) {
	ing := obj.(*v1beta1.Ingress)
	rules := ing.Spec.Rules
	hosts := make([]string, len(rules))
	for index, rule := range rules {
		log.Debugf("create %s\n", rule.Host)
		hosts[index] = i.getDomainName(rule.Host)
	}
	if change := i.addData(hosts...); change {
		i.updateEtcd(i.data.Elements())
	}
}

func (i *IngressWatcher) del(obj interface{}) {
	ing := obj.(*v1beta1.Ingress)
	rules := ing.Spec.Rules
	hosts := make([]string, len(rules))
	for index, rule := range rules {
		log.Debugf("delete %s\n", rule.Host)
		hosts[index] = i.getDomainName(rule.Host)
	}
	if change := i.deleteData(hosts...); change {
		i.updateEtcd(i.data.Elements())
	}
}

func (i *IngressWatcher) update(old, cur interface{}) {
	oldIng := old.(*v1beta1.Ingress)
	curIng := cur.(*v1beta1.Ingress)
	oldRules := oldIng.Spec.Rules
	curRules := curIng.Spec.Rules
	delHosts := make([]string, len(oldRules))
	addHosts := make([]string, len(curRules))
	for index, rule := range curRules {
		addHosts[index] = i.getDomainName(rule.Host)
	}
	for index, rule := range oldRules {
		delHosts[index] = i.getDomainName(rule.Host)
	}
	for index, host := range delHosts {
		for _, host2 := range addHosts{
			if host == host2 {
				delHosts[index] = ""
			}
		}
	}
	delChange := i.deleteData(delHosts...)
	addChange := i.addData(addHosts...)
	if delChange || addChange {
		i.updateEtcd(i.data.Elements())
	}
}
