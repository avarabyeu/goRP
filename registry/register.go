package registry

import (
	"log"
	"time"
)

const (
	RETRY_TIMEOUT  time.Duration = time.Second * 5
	RETRY_ATTEMPTS int           = 3
)

type ServiceDiscovery interface {
	Register() error
	Deregister() error
}

func Register(discovery ServiceDiscovery) {
	err := tryRegister(discovery)
	if nil != err {
		log.Fatal(err)
	}

	shutdownHook(func() error {
		log.Println("try to deregister")
		return Deregister(discovery)

	})
}

func Deregister(discovery ServiceDiscovery) error {
	return tryDeregister(discovery)
}

func tryRegister(discovery ServiceDiscovery) error {
	return retry(RETRY_ATTEMPTS, RETRY_TIMEOUT, func() error {
		e := discovery.Register()
		if nil != e {
			log.Printf("Cannot register service: %s", e)
		} else {
			log.Print("Service Registered!")
		}
		return e
	})
}

func tryDeregister(discovery ServiceDiscovery) error {
	return retry(RETRY_ATTEMPTS, RETRY_TIMEOUT, func() error {
		return discovery.Deregister()
	})
}
