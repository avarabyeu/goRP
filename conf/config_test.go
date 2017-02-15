package conf

import (
	"testing"
	"os"
)

func TestLoadConfig(t *testing.T) {
	rpConf := LoadConfig("./../server.yaml", nil)
	if "10.200.10.1" != rpConf.Server.Hostname {
		t.Error("Config parser fails")
	}
}

func TestLoadConfigWithParameters(t *testing.T) {
	os.Setenv("RP_OK", "param1")
	os.Setenv("RP_PARAMETERS.PARAM", "env_value")
	rpConf := LoadConfig("", map[string]interface{}{"parameters.param": "default_value"})

	if "env_value" != rpConf.Param("parameters.param") {
		t.Error("Config parser fails")
	}
}

func TestLoadConfigUnexisted(t *testing.T) {
	rpConf := LoadConfig("server.yaml", nil)
	if "" != rpConf.Server.Hostname {
		t.Error("Should return empty string for default config")
	}
}

func TestLoadConfigIncorrectFormat(t *testing.T) {
	rpConf := LoadConfig("config_test.go", nil)
	if "" != rpConf.Server.Hostname {
		t.Error("Should return empty string for default config")
	}
}
