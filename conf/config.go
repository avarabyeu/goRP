package conf

import (
	"github.com/spf13/viper"
	"log"
)

type ServerConfig struct {
	Hostname string
	Port     int
}

type EurekaConfig struct {
	Url          string
	AppName      string
	PollInterval int
}

type ConsulConfig struct {
	Address      string
	Scheme       string
	Token        string
	AppName      string
	PollInterval int
	Tags         []string
}

type RpConfig struct {
	Server ServerConfig
	Eureka EurekaConfig
	Consul ConsulConfig
}

func LoadConfig(file string) *RpConfig {
	var config = viper.New();

	config.SetConfigType("yaml")
	config.SetConfigFile(file)
	config.AutomaticEnv()
	err := config.ReadInConfig()
	if err != nil {
		log.Fatal("No configuration file loaded - using defaults")
	}

	var rpConf RpConfig
	config.Unmarshal(&rpConf)
	return &rpConf
}
