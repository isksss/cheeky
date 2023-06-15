package config

import "os/exec"

// Gitがインストールされているかどうかを確認する
func IsGitInstalled() bool {
	_, err := exec.LookPath("git")
	if err != nil {
		return false
	}
	return true
}
