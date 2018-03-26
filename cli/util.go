package cli

import (
	"errors"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

func validateConfig(cfg *config) error {
	if "" == cfg.UUID {
		return errors.New("uuid is not set")
	}

	if "" == cfg.Project {
		return errors.New("project is not set")
	}

	if "" == cfg.Host {
		return errors.New("host is not set")
	}
	return nil
}

func answerYes(answer string) bool {
	lower := strings.ToLower(answer)
	return "y" == lower || "yes" == lower
}

func configFilePresent() bool {
	_, err := os.Stat(getConfigFile())
	return !os.IsNotExist(err)
}

func getConfigFile() string {
	return filepath.Join(getHomeDir(), ".gorp")
}
func getHomeDir() string {
	if h := os.Getenv("HOME"); "" != h {
		return h
	}
	curUser, err := user.Current()
	if err != nil {
		// well, sheesh
		return "."
	}

	return curUser.HomeDir
}
