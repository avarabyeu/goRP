package conf

import "testing"

func TestLoadConfig(t *testing.T) {
	rpConf := LoadConfig("./../server.yaml")
	if "10.200.10.1" != rpConf.Server.Hostname {
		t.Error("Config parser fails")
	}
}

func TestLoadConfigUnexisted(t *testing.T) {
	rpConf := LoadConfig("server.yaml")
	if rpConf.Server.Hostname != "" {
		t.Error("Should return empty string for default config")
	}
}

func TestLoadConfigIncorrectFormat(t *testing.T) {
	rpConf := LoadConfig("config_test.go")
	if rpConf.Server.Hostname != "" {
		t.Error("Should return empty string for default config")
	}
}
