package conf

import (
	"log"

	"fmt"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"path/filepath"
	"strings"
	"os"
)

//Registry represents type of used service discovery server
type Registry string

const (
	//Consul service discovery
	Consul Registry = "consul"
	//Eureka service discovery
	Eureka Registry = "eureka"
)

//ServerConfig represents Main service configuration
type ServerConfig struct {
	Hostname string
	Port     int
}

//EurekaConfig represents Eureka Discovery service configuration
type EurekaConfig struct {
	URL          string
	PollInterval int
}

//ConsulConfig represents Consul Discovery service configuration
type ConsulConfig struct {
	Address      string
	Scheme       string
	Token        string
	PollInterval int
	Tags         string
}

//RpConfig represents Composite of all app configs
type RpConfig struct {
	AppName  string
	Registry Registry
	Server   ServerConfig
	Eureka   EurekaConfig
	Consul   ConsulConfig

	rawConfig *viper.Viper
}

//Get reads parameter/property value from config (env,file,defaults)
func (cfg *RpConfig) Get(key string) interface{} {
	return cfg.rawConfig.Get(key)
}

//LoadConfig loads configuration from provided file and serializes it into RpConfig struct
func LoadConfig(file string, defaults map[string]interface{}) *RpConfig {
	var vpr = viper.New()

	if "" != file {
		vpr.SetConfigType(strings.TrimLeft(filepath.Ext(file), "."))
		vpr.SetConfigFile(file)
	}

	vpr.SetEnvPrefix("RP")
	vpr.AutomaticEnv()

	applyDefaults(vpr)
	if nil != defaults {
		for k, v := range defaults {
			vpr.SetDefault(k, v)
		}
	}

	err := vpr.ReadInConfig()
	if err != nil {
		log.Println("No configuration file loaded - using defaults")
	}

	bindToFlags(vpr)

	var rpConf RpConfig
	err = vpr.Unmarshal(&rpConf)
	if err != nil {
		log.Fatalf("Cannot unmarshal config: %s", err.Error())
	}
	rpConf.rawConfig = vpr

	//vpr.Debug()
	return &rpConf
}

func bindToFlags(vpr *viper.Viper) {
	if !pflag.Parsed() {
		for _, key := range vpr.AllKeys() {
			pflag.String(key, vpr.GetString(key), fmt.Sprintf("Property: %s", key))
		}

		pflag.Parse()

		pflag.VisitAll(func(flg *pflag.Flag) {
			if "" != flg.Value.String() {
				vpr.BindPFlag(flg.Name, flg)
			}
		})

	}

}

func applyDefaults(vpr *viper.Viper) {

	vpr.SetDefault("appname", "goRP")

	vpr.SetDefault("registry", Consul)

	vpr.SetDefault("server.port", 8080)

	defaultHostname := os.Getenv("HOSTNAME")
	if "" == defaultHostname {
		defaultHostname = "localhost"
	}
	vpr.SetDefault("server.hostname", defaultHostname)

	vpr.SetDefault("eureka.url", "http://localhost:8761/eureka")
	vpr.SetDefault("eureka.appname", "goRP")
	vpr.SetDefault("eureka.pollInterval", 5)

	vpr.SetDefault("consul.address", "localhost:8500")
	vpr.SetDefault("consul.scheme", "http")
	vpr.SetDefault("consul.pollInterval", 5)
	vpr.SetDefault("consul.tags", "")

}
