package registry

import (
	"log"
	"github.com/avarabyeu/goRP/conf"
	"github.com/hashicorp/consul/api"
	"strconv"
	"fmt"
)

type ConsulClient struct {
	c   *api.Client
	reg *api.AgentServiceRegistration
}

func NewConsul(cfg *conf.RpConfig) ServiceDiscovery {
	c, err := api.NewClient(&api.Config{
		Address: cfg.Consul.Address,
		Scheme: cfg.Consul.Scheme,
		Token: cfg.Consul.Token})
	if nil != err {
		log.Fatal("Cannot create Consul client!")
	}

	baseUrl := PROTOCOL + cfg.Server.Hostname + ":" + strconv.Itoa(cfg.Server.Port)
	registration := &api.AgentServiceRegistration{
		ID :     fmt.Sprintf("%s-%s-%d", cfg.Consul.AppName, cfg.Server.Hostname, cfg.Server.Port),
		Port: cfg.Server.Port,
		Address: getLocalIP(),
		Name: cfg.Consul.AppName,
		Tags: cfg.Consul.Tags,
		Check: &api.AgentServiceCheck{
			HTTP: baseUrl + "/health",
			Interval: fmt.Sprintf("%ds", cfg.Consul.PollInterval),
		},

	}
	return &ConsulClient{
		c: c,
		reg:registration,
	};

}

func (ec *ConsulClient) Register() error {
	return ec.c.Agent().ServiceRegister(ec.reg)
}

func (ec *ConsulClient) Deregister() error {
	e := ec.c.Agent().ServiceDeregister(ec.reg.ID)
	if nil != e {
		log.Print(e)
	}
	return e
}

