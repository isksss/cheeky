package config

import (
	"os"
	"path/filepath"
)

func GetSSHDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".ssh")
}

// fileを削除する
func RemoveFile(file string) error {
	return os.Remove(file)
}

// シンボリックリンクを作成する
func CreateSymlink(src string, dst string) error {
	return os.Symlink(src, dst)
}
