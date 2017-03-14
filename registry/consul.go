package registry

import (
	"fmt"
	"log"
	"strconv"

	"github.com/avarabyeu/goRP/commons"
	"github.com/avarabyeu/goRP/conf"
	"github.com/hashicorp/consul/api"
	"strings"
)

type consulClient struct {
	c   *api.Client
	reg *api.AgentServiceRegistration
}

//NewConsul creates new instance of Consul implementation of ServiceDiscovery
//based on provided config
func NewConsul(cfg *conf.RpConfig) ServiceDiscovery {
	c, err := api.NewClient(&api.Config{
		Address: cfg.Consul.Address,
		Scheme:  cfg.Consul.Scheme,
		Token:   cfg.Consul.Token})
	if nil != err {
		log.Fatal("Cannot create Consul client!")
	}

	baseURL := commons.HTTP + cfg.Server.Hostname + ":" + strconv.Itoa(cfg.Server.Port)
	registration := &api.AgentServiceRegistration{
		ID:      fmt.Sprintf("%s-%s-%d", cfg.AppName, cfg.Server.Hostname, cfg.Server.Port),
		Port:    cfg.Server.Port,
		Address: commons.GetLocalIP(),
		Name:    cfg.AppName,
		Tags:    strings.Split(cfg.Consul.Tags, ","),
		Check: &api.AgentServiceCheck{
			HTTP:     baseURL + "/health",
			Interval: fmt.Sprintf("%ds", cfg.Consul.PollInterval),
		},
	}
	return &consulClient{
		c:   c,
		reg: registration,
	}

}

//Register registers instance in Consul
func (ec *consulClient) Register() error {
	return ec.c.Agent().ServiceRegister(ec.reg)
}

//Deregister de-registers instance in Consul
func (ec *consulClient) Deregister() error {
	e := ec.c.Agent().ServiceDeregister(ec.reg.ID)
	if nil != e {
		log.Print(e)
	}
	return e
}

//DoWithClient does provided action using service discovery client
func (ec *consulClient) DoWithClient(f func(client interface{}) (interface{}, error)) (interface{}, error) {
	return f(ec.c)
}
