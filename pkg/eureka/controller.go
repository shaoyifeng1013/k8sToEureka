package eureka

import (
	"github.com/ArthurHlt/go-eureka-client/eureka"
	"github.com/robfig/cron/v3"
	"k8s.io/klog/v2"
)

// pre 5 second
const spec = "@every 5s"

type Controller struct {
	clients           map[string]*eureka.Client
	instances         map[string]*eureka.InstanceInfo
	cornJobIds        map[string]cron.EntryID
	cornClient        *cron.Cron
	DefaultEurekaAddr string
}

func (e *Controller) InitController() {
	e.clients = make(map[string]*eureka.Client, 8)
	e.cornJobIds = make(map[string]cron.EntryID, 8)
	e.cornClient = cron.New()
	e.DefaultEurekaAddr = "http://eureka.didispace.com/eureka"
}

//TODO corn job 少一个优雅退出
//TODO 重试机制
//TODO eureka 创建加个限流
func (e *Controller) CreateClient(eurekaAddr, hostName, app, ip string, port int) error {
	if eurekaAddr == "" {
		eurekaAddr = e.DefaultEurekaAddr
	}
	client := eureka.NewClient([]string{
		eurekaAddr,
	})
	instance := eureka.NewInstanceInfo(hostName, app, ip, port, 80, false)
	instance.Metadata = &eureka.MetaData{
		Map: map[string]string{"sourceBy": "syfclient"},
	}
	err := client.RegisterInstance(app, instance)
	if err != nil {
		return err
	}
	e.clients[app] = client
	jobId, err := e.cornClient.AddJob(spec, &eurekaHealthCheck{
		client:   client,
		instance: instance,
	})
	if err != nil {
		return err
	}
	e.cornJobIds[instance.IpAddr] = jobId
	e.cornClient.Start()
	klog.Info("service ", app, " ", ip, " register success")
	return nil
}

func (e *Controller) Offline(app, ip string) error {
	client := e.clients[app]
	err := client.UnregisterInstance(app, ip)
	if err != nil {
		return err
	}
	//删心跳
	e.cornClient.Remove(e.cornJobIds[ip])
	klog.Info("service ", app, " ", ip, " offline success")
	return nil
}

func (e *Controller) Exit() {
	if e.cornClient != nil {
		e.cornClient.Stop()
	}
}

type eurekaHealthCheck struct {
	client   *eureka.Client
	instance *eureka.InstanceInfo
}

func (e *eurekaHealthCheck) Run() {
	err := e.client.SendHeartbeat(e.instance.App, e.instance.HostName)
	if err == nil {
		klog.Info(e.instance.App, " ", e.instance.IpAddr+" heartbeat")
	} else {
		klog.Info(e.instance.App, " ", e.instance.IpAddr+" heartbeat send error "+err.Error())
	}
}
