package conf

import (
	"os"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	rpConf := LoadConfig("./../server.yaml", nil)
	if "10.200.10.1" != rpConf.Server.Hostname {
		t.Error("Config parser fails")
	}
}

func TestLoadConfigWithParameters(t *testing.T) {
	os.Setenv("RP_PARAMETERS.PARAM", "env_value")
	rpConf := LoadConfig("", map[string]interface{}{"parameters.param": "default_value"})

	if "env_value" != rpConf.Get("parameters.param").(string) {
		t.Error("Config parser fails")
	}
}

func TestLoadConfigNonExisting(t *testing.T) {
	rpConf := LoadConfig("server.yaml", nil)
	if 8080 != rpConf.Server.Port {
		t.Error("Should return empty string for default config")
	}
}

func TestLoadConfigIncorrectFormat(t *testing.T) {
	rpConf := LoadConfig("config_test.go", nil)
	if 8080 != rpConf.Server.Port {
		t.Error("Should return empty string for default config")
	}
}
