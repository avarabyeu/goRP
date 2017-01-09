package reportportal

import (
	"time"
	"log"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"net"
)

func retry(attempts int, timeout time.Duration, callback func() error) (err error) {
	for i := 0; i <= attempts - 1; i++ {
		err = callback()
		if err == nil {
			return nil
		}

		//time.Sleep(timeout)
		<-time.After(timeout)
		log.Println("retrying...")
	}
	return fmt.Errorf("after %d attempts, last error: %s", attempts, err)
}

func shutdownHook(hook func() error) {
	c := make(chan os.Signal, 1)          // Create a channel accepting os.Signal
	// Bind a given os.Signal to the channel we just created
	signal.Notify(c, os.Interrupt)        // Register os.Interrupt
	signal.Notify(c, syscall.SIGTERM)     // Register syscall.SIGTERM

	go func() {
		// Start an anonymous func running in a goroutine
		<-c                           // that will block until a message is recieved on
		e := hook()
		if nil != e {
			log.Println("Shutdown hook error: ", e)
		}

		os.Exit(1)
	}()
}

func getLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}
