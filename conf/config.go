package conf

import (
	"log"

	"github.com/spf13/viper"
)

//ServerConfig represents Main service configuration
type ServerConfig struct {
	Hostname string
	Port     int
}

//EurekaConfig represents Eureka Discovery service configuration
type EurekaConfig struct {
	URL          string
	AppName      string
	PollInterval int
}

//ConsulConfig represents Consul Discovery service configuration
type ConsulConfig struct {
	Address      string
	Scheme       string
	Token        string
	AppName      string
	PollInterval int
	Tags         []string
}

//RpConfig represents Composite of all app configs
type RpConfig struct {
	Server ServerConfig
	Eureka EurekaConfig
	Consul ConsulConfig
}

//LoadConfig loads configuration from provided file and serializes it into RpConfig struct
func LoadConfig(file string) *RpConfig {
	var config = viper.New()

	config.SetConfigType("yaml")
	config.SetConfigFile(file)
	config.AutomaticEnv()
	err := config.ReadInConfig()
	if err != nil {
		log.Println("No configuration file loaded - using defaults")
	}

	var rpConf RpConfig
	config.Unmarshal(&rpConf)
	return &rpConf
}
