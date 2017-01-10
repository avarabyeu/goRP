package registry

import (
	"github.com/avarabyeu/goRP/conf"
	"github.com/hudl/fargo"
	"log"
	"strconv"
	"time"
)

const PROTOCOL = "http://"

type EurekaClient struct {
	eureka      fargo.EurekaConnection
	appInstance *fargo.Instance
}

func NewEureka(conf *conf.RpConfig) ServiceDiscovery {
	eureka := fargo.NewConn(conf.Eureka.Url)
	eureka.PollInterval = time.Duration(conf.Eureka.PollInterval) * time.Second
	baseUrl := PROTOCOL + conf.Server.Hostname + ":" + strconv.Itoa(conf.Server.Port)
	var appInstance = &fargo.Instance{
		App:        conf.Eureka.AppName,
		VipAddress: conf.Server.Hostname,
		IPAddr:     getLocalIP(),
		HostName:   conf.Server.Hostname,
		Port:       conf.Server.Port,
		DataCenterInfo: fargo.DataCenterInfo{
			Name: "MyOwn",
		},
		HomePageUrl:    baseUrl + "/",
		HealthCheckUrl: baseUrl + "/health",
		StatusPageUrl:  baseUrl + "/info",
		Status:         fargo.UP,
	}
	ec := &EurekaClient{
		eureka:      eureka,
		appInstance: appInstance,
	}
	return ec
}

func (ec *EurekaClient) Register() error {
	e := ec.eureka.RegisterInstance(ec.appInstance)
	if nil == e {
		heartBeat(ec)
	}
	return e
}

func (ec *EurekaClient) Deregister() error {
	return ec.eureka.DeregisterInstance(ec.appInstance)
}

func heartBeat(ec *EurekaClient) {
	go func() {
		for {
			err := ec.eureka.HeartBeatInstance(ec.appInstance)
			if err != nil {
				code, ok := fargo.HTTPResponseStatusCode(err)
				if ok && 404 == code {
					Register(ec)
				}
				log.Printf("Failure updating %s in goroutine", ec.appInstance.Id())
			}
			<-time.After(ec.eureka.PollInterval)
		}
	}()
}
