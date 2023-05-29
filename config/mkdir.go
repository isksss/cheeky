package config

import (
	"errors"
	"os"
	"path/filepath"
)

func MakeDir() {
	CreateDirIfNotExist(GetCheekyHome())
	CreateDirIfNotExist(filepath.Join(GetCheekyHome(), "keys"))
}

func CreateDirIfNotExist(dir string) error {
	if _, err := os.Stat(dir); !os.IsNotExist(err) {
		return errors.New("directory already exists")
	}
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	return nil
}
