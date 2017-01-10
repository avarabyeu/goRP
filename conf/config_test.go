package conf

import (
	"testing"
)

func TestLoadConfig(t *testing.T) {

	var rpConf = LoadConfig("./../server.yaml")
	if "10.200.10.1" != rpConf.Server.Hostname {
		t.Error("Config parser fails");
	}
}
