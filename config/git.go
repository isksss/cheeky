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

type GitConfig struct {
	Key   string
	Value string
}

func SetGitConfig(config GitConfig) error {
	cmd := exec.Command("git", "config", config.Key, config.Value)
	return cmd.Run()
}

func GetGitConfig(config GitConfig) (string, error) {
	cmd := exec.Command("git", "config", config.Key)
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(out), nil
}

func SetGitUser(name string) error {
	config := GitConfig{
		Key:   "user.name",
		Value: name,
	}

	cmd := exec.Command("git", "config", "--global", config.Key, config.Value)
	return cmd.Run()
}

func SetGitEmail(email string) error {
	config := GitConfig{
		Key:   "user.email",
		Value: email,
	}

	cmd := exec.Command("git", "config", "--global", config.Key, config.Value)
	return cmd.Run()
}
