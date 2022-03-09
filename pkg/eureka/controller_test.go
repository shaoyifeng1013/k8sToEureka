package eureka

import (
	"fmt"
	"github.com/robfig/cron/v3"
	"testing"
	"time"
)

func TestCreateClient(t *testing.T) {
	controller := Controller{}
	controller.InitController()
	err := controller.CreateClient("http://eureka.didispace.com/eureka", "127.0.0.2", "syf-test", "127.0.0.1", 80)
	if err != nil {
		t.Error(err)
	}

	app := make(chan string)
	go func(app chan string) {
		time.Sleep(5 * time.Second)
		client := controller.clients["syf-test"]
		application, err := client.GetApplication("syf-test")
		if err != nil {
			app <- "error " + err.Error()
		}
		app <- application.Instances[0].IpAddr
	}(app)

	select {
	case instance := <-app:
		fmt.Println(instance)
	case <-time.After(2 * time.Second):
		fmt.Println("timeout")
	}
	for {

	}
}
func TestCron(t *testing.T) {
	t.Log(cron.Every(5 * time.Second).Delay.String())
}
