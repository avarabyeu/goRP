package conf

import (
	"testing"
)

func TestLoadConfig(t *testing.T) {

	var rpConf = LoadConfig("./../server.yaml")
	if "localhost" != rpConf.Server.Hostname {
		t.Error("Config parser fails");
	}
}
