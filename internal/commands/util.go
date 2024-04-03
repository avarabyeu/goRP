package commands

import (
	"errors"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

var (
	errHostNotSet    = errors.New("host is not set")
	errProjectNotSet = errors.New("project is not set")
	errUUIDNotSet    = errors.New("uuid is not set")
)

func validateConfig(cfg *config) error {
	if cfg.UUID == "" {
		return errUUIDNotSet
	}

	if cfg.Project == "" {
		return errProjectNotSet
	}

	if cfg.Host == "" {
		return errHostNotSet
	}

	return nil
}

func answerYes(answer string) bool {
	lower := strings.ToLower(answer)

	return lower == "y" || lower == "yes"
}

func configFilePresent() bool {
	_, err := os.Stat(getConfigFile())

	return !os.IsNotExist(err)
}

func getConfigFile() string {
	return filepath.Join(getHomeDir(), ".gorp")
}

func getHomeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	curUser, err := user.Current()
	if err != nil {
		// well, sheesh
		return "."
	}

	return curUser.HomeDir
}
