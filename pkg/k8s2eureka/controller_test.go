package k8s2eureka

import (
	"context"
	"fmt"
	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/tools/cache"
	"k8s.io/klog/v2"
	"reflect"
	"sigs.k8s.io/yaml"
	"study/pkg/apis/runtime/v1alpha"
	"study/pkg/config"
	"study/pkg/crd"
	"sync"
	"testing"
)

func TestControllerDeployList(t *testing.T) {
	conf := &config.Config{
		Name:           "test",
		KubeConfigPath: `D:/ForCoding/zybank/k8sToEureka/config`,
	}

	controller := CreateController()
	e := controller.InitController(*conf)
	if e != nil {
		panic(e)
	}

	deploymentList, e := controller.DefaultClientSet.AppsV1().Deployments("istio-system").List(context.TODO(), v1.ListOptions{})
	if e != nil {
		panic(e)
	}
	for _, deploy := range deploymentList.Items {
		fmt.Println(deploy.Name)
	}
}

func TestCreateDeploy(t *testing.T) {
	conf := &config.Config{
		Name:           "test",
		KubeConfigPath: `D:/ForCoding/zybank/k8sToEureka/config`,
	}

	controller := CreateController()
	e := controller.InitController(*conf)
	if e != nil {
		panic(e)
	}
	deploymentClient := controller.DefaultClientSet.AppsV1().Deployments(v1.NamespaceDefault)
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: "demo-deployment",
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: int32Ptr(2),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "demo",
				},
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "demo",
					},
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:  "web",
							Image: "nginx:1.12",
							Ports: []apiv1.ContainerPort{
								{
									Name:          "http",
									Protocol:      apiv1.ProtocolTCP,
									ContainerPort: 80,
								},
							},
						},
					},
				},
			},
		},
	}

	// Create Deployment
	fmt.Println("Creating deployment...")
	result, err := deploymentClient.Create(context.TODO(), deployment, metav1.CreateOptions{})
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
}

func int32Ptr(i int32) *int32 {
	return &i
}

func TestUpdateDeploy(t *testing.T) {
	conf := &config.Config{
		Name:           "test",
		KubeConfigPath: `D:/ForCoding/zybank/k8sToEureka/config`,
	}

	controller := CreateController()
	e := controller.InitController(*conf)
	if e != nil {
		panic(e)
	}
	deploymentClient := controller.DefaultClientSet.AppsV1().Deployments(v1.NamespaceDefault)
	deployment, e := deploymentClient.Get(context.TODO(), "demo-deployment", v1.GetOptions{})
	if e != nil {
		panic(e)
	}
	deployment.Spec.Replicas = int32Ptr(1)
	fmt.Println("Updating deployment...")
	result, e := deploymentClient.Update(context.TODO(), deployment, v1.UpdateOptions{})
	fmt.Println(result, e)

	//delete
	fmt.Println("Deleting deployment...")
	err := deploymentClient.Delete(context.TODO(), "demo-deployment", v1.DeleteOptions{})
	fmt.Println(err)
}

func TestWatchPods(t *testing.T) {
	conf := &config.Config{
		Name:           "test",
		KubeConfigPath: `D:/ForCoding/zybank/k8sToEureka/config`,
	}

	controller := CreateController()
	e := controller.InitController(*conf)
	if e != nil {
		panic(e)
	}
	stop := make(chan struct{}, 0)
	var lock sync.WaitGroup
	lock.Add(1)
	podListWatcher := cache.NewListWatchFromClient(controller.DefaultClientSet.CoreV1().RESTClient(), "pods", v1.NamespaceDefault, fields.Everything())
	funcs := cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			klog.Info("AddFunc\n", obj)
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			klog.Info("UpdateFunc\n", oldObj, newObj)
		},
		DeleteFunc: func(obj interface{}) {
			klog.Info("DeleteFunc\n", obj)
			lock.Done()
		},
	}
	_, k8sController := cache.NewInformer(podListWatcher, &apiv1.Pod{}, 0, funcs)

	//run
	go k8sController.Run(stop)

	lock.Wait()
}

func TestCreateCrd(t *testing.T) {
	conf := &config.Config{
		Name:           "test",
		KubeConfigPath: `D:/ForCoding/zybank/k8sToEureka/config`,
	}

	controller := CreateController()
	e := controller.InitController(*conf)
	if e != nil {
		panic(e)
	}
	u := &unstructured.Unstructured{}
	u.SetAPIVersion(crd.MeshServers.ApiVersion)
	u.SetKind(crd.MeshServers.Kind)
	u.SetNamespace("default")
	u.SetName("test-server2")
	u.SetLabels(map[string]string{
		"version": "v1",
	})
	u.SetAnnotations(map[string]string{
		"user": "syf",
	})

	meshServer := &crd.MeshServer{
		Server: crd.Server{
			Host: "test-server2.com",
			Port: 8080,
		},
		Zone: "cluster1",
	}
	meshServerBytes, err := yaml.Marshal(meshServer)
	if err != nil {
		panic(e)
	}
	ms := &map[string]interface{}{}
	e = yaml.Unmarshal(meshServerBytes, ms)
	if err != nil {
		panic(e)
	}
	u.Object["spec"] = ms
	result, err := controller.DynamicSet.Resource(crd.MeshServers.GroupVersionResource).Namespace("default").Create(context.TODO(), u, metav1.CreateOptions{})
	fmt.Println(result, err)
}

func TestInformer(t *testing.T) {
	conf := &config.Config{
		Name:           "test",
		KubeConfigPath: `D:/ForCoding/zybank/k8sToEureka/config`,
	}

	stop := make(chan struct{})

	meshservercontroller := CreateController()
	e := meshservercontroller.InitController(*conf)
	if e != nil {
		panic(e)
	}
	meshservercontroller.Run(stop)
}

func TestDeepEqual(t *testing.T) {
	server1 := v1alpha.MeshServer{
		TypeMeta:   metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{},
		Spec: v1alpha.MeshServerSpec{
			Server: v1alpha.Server{
				Host:   "1",
				Port:   0,
				IpAddr: "1",
			},
			Zone: "1",
		},
	}

	server2 := v1alpha.MeshServer{
		TypeMeta:   metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{},
		Spec: v1alpha.MeshServerSpec{
			Server: v1alpha.Server{
				Host:   "1",
				Port:   0,
				IpAddr: "1",
			},
			Zone: "1",
		},
	}

	t.Log(reflect.DeepEqual(server2, server1))
}
