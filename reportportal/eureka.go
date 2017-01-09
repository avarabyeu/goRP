package reportportal

import (
	"github.com/hudl/fargo"
	"log"
	"time"
	"strconv"
)

func registerInEureka(rp *RpServer) {
	eureka := fargo.NewConn(rp.conf.Eureka.Url)

	eureka.PollInterval = POLL_INTERVAL

	baseUrl := PROTOCOL + rp.conf.Server.Hostname + ":" + strconv.Itoa(rp.conf.Server.Port)

	var appInstance = &fargo.Instance{
		App: rp.conf.Eureka.AppName,
		VipAddress: rp.conf.Server.Hostname,
		IPAddr: getLocalIP(),
		HostName: rp.conf.Server.Hostname,
		Port: rp.conf.Server.Port,
		DataCenterInfo: fargo.DataCenterInfo{
			Name:"MyOwn",
		},
		HomePageUrl:baseUrl + "/",
		HealthCheckUrl: baseUrl + "/health",
		StatusPageUrl:baseUrl + "/info",
		Status: fargo.UP,

	}

	err := tryRegister(&eureka, appInstance)
	if nil != err {
		log.Fatal(err)
	}
	heartBeat(&eureka, appInstance)

	shutdownHook(func() error {
		return retry(RETRY_ATTEMPTS, RETRY_TIMEOUT, func() error {
			return eureka.DeregisterInstance(appInstance)
		})

	})

}

func tryRegister(eureka *fargo.EurekaConnection, instance *fargo.Instance) error {
	return retry(RETRY_ATTEMPTS, RETRY_TIMEOUT, func() error {
		return eureka.RegisterInstance(instance)
	})

}

func heartBeat(eureka *fargo.EurekaConnection, instance *fargo.Instance) {
	go func() {
		for {
			err := eureka.HeartBeatInstance(instance)
			if err != nil {
				code, ok := fargo.HTTPResponseStatusCode(err)
				if ok && 404 == code {
					tryRegister(eureka, instance)
				}
				log.Printf("Failure updating %s in goroutine", instance.Id())
			}
			<-time.After(eureka.PollInterval)
		}
	}()
}
