package config

import (
	"os"
	"path/filepath"
)

func GetCheekyHome() string {
	return filepath.Join(getXdgConfigHome(), "cheeky")
}

func getXdgConfigHome() string {
	xdg := os.Getenv("XDG_CONFIG_HOME")
	if xdg != "" {
		return xdg
	}

	dir, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(dir, ".config")
}

func GetKeysDir() string {
	return filepath.Join(GetCheekyHome(), "keys")
}
