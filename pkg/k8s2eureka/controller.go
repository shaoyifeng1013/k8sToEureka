package k8s2eureka

import (
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/klog/v2"
	"reflect"
	"study/pkg/apis/runtime/v1alpha"
	"study/pkg/config"
	"study/pkg/eureka"
	"study/pkg/generated/clientset/versioned"
	"study/pkg/generated/informers/externalversions"
	"time"
)
import "k8s.io/client-go/tools/clientcmd"

type Controller struct {
	DefaultClientSet    *kubernetes.Clientset
	DynamicSet          dynamic.Interface
	MeshServerClientSet *versioned.Clientset
	EurekaController    *eureka.Controller
}

func (k *Controller) InitController(config config.Config) error {
	/**
	restConfig not nil 则继续， nil 则incluster
	*/
	var restConfig *rest.Config
	var err error
	if config.KubeConfigPath != "" {
		restConfig, err = clientcmd.BuildConfigFromFlags("", config.KubeConfigPath)
		if err != nil {
			klog.Info("KubeConfigPath load error")
			return err
		}
	} else {
		restConfig, err = rest.InClusterConfig()
		if err != nil {
			klog.Info("InClusterConfig load error")
			return err
		}
	}

	cs, e := kubernetes.NewForConfig(restConfig)
	if e != nil {
		return e
	}
	k.DefaultClientSet = cs

	dynamicset, e := dynamic.NewForConfig(restConfig)
	if e != nil {
		return e
	}
	k.DynamicSet = dynamicset

	MeshClientSet, e := versioned.NewForConfig(restConfig)
	if e != nil {
		return e
	}
	k.MeshServerClientSet = MeshClientSet

	eurekaController := &eureka.Controller{}
	eurekaController.InitController()
	k.EurekaController = eurekaController
	return nil
}

func CreateController() *Controller {
	kube := &Controller{}
	return kube
}

func (c Controller) Run(stopCh <-chan struct{}) error {
	defer func() {
		c.EurekaController.Exit()
	}()
	//创建informer
	factory := externalversions.NewSharedInformerFactory(c.MeshServerClientSet, 3*time.Second)
	informer := factory.Runtime().V1alpha().MeshServers()
	informer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			meshServer := obj.(*v1alpha.MeshServer)
			klog.Info("AddFunc ", meshServer.Name)
			if !meshServer.IsValid() {
				return
			}
			err := c.EurekaController.CreateClient("http://localhost:8761/eureka", meshServer.Spec.IpAddr, meshServer.Spec.Host, meshServer.Spec.IpAddr, meshServer.Spec.Port)
			if err != nil {
				klog.Error(err)
			}
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			klog.Info("UpdateFunc")
			if reflect.DeepEqual(oldObj, newObj) {
				klog.Info("same change")
				return
			}
			newMeshServer := newObj.(*v1alpha.MeshServer)
			oldMeshServer := oldObj.(*v1alpha.MeshServer)

			klog.Info(newMeshServer, "\n", oldMeshServer)
		},
		DeleteFunc: func(obj interface{}) {
			klog.Info("DeleteFunc")
			meshServer := obj.(*v1alpha.MeshServer)
			if !meshServer.IsValid() {
				return
			}
			err := c.EurekaController.Offline(meshServer.Spec.Host, meshServer.Spec.IpAddr)
			if err != nil {
				klog.Warning(err)
			}
		},
	})

	factory.Start(stopCh)
	klog.Info("Started Kube Controller")
	<-stopCh
	klog.Info("Shutting down Kube Controller")
	return nil
}
