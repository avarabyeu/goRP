package registry

import (
	"github.com/avarabyeu/goRP/conf"
	"github.com/hudl/fargo"
	"log"
	"strconv"
	"time"
	"github.com/avarabyeu/goRP/common"
)

type eurekaClient struct {
	eureka      fargo.EurekaConnection
	appInstance *fargo.Instance
}

//NewEureka creates new instance of Eureka implementation of Service Discovery
func NewEureka(conf *conf.RpConfig) ServiceDiscovery {
	eureka := fargo.NewConn(conf.Eureka.URL)
	eureka.PollInterval = time.Duration(conf.Eureka.PollInterval) * time.Second
	baseURL := common.HTTP + conf.Server.Hostname + ":" + strconv.Itoa(conf.Server.Port)
	var appInstance = &fargo.Instance{
		App:        conf.AppName,
		VipAddress: conf.Server.Hostname,
		IPAddr:     common.GetLocalIP(),
		HostName:   conf.Server.Hostname,
		Port:       conf.Server.Port,
		DataCenterInfo: fargo.DataCenterInfo{
			Name: "MyOwn",
		},
		HomePageUrl:    baseURL + "/",
		HealthCheckUrl: baseURL + "/health",
		StatusPageUrl:  baseURL + "/info",
		Status:         fargo.UP,
	}
	ec := &eurekaClient{
		eureka:      eureka,
		appInstance: appInstance,
	}
	return ec
}

//Register registers instance in Eureka
func (ec *eurekaClient) Register() error {
	e := ec.eureka.RegisterInstance(ec.appInstance)
	if nil == e {
		heartBeat(ec)
	}
	return e
}

//Deregister de-registers instance in Eureka
func (ec *eurekaClient) Deregister() error {
	return ec.eureka.DeregisterInstance(ec.appInstance)
}

//DoWithClient does provided action using service discovery client
func (ec *eurekaClient) DoWithClient(f func(client interface{}) (interface{}, error)) (interface{}, error) {
	return f(ec.eureka)
}

//sends heartbeats to Eureka to notify it that service is still alive
func heartBeat(ec *eurekaClient) {
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
