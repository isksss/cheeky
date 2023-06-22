package commands

import (
	"context"
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/google/subcommands"
	"github.com/isksss/cheeky/config"
)

type SetCmd struct{}

func (*SetCmd) Name() string {
	return "set"
}

func (*SetCmd) Synopsis() string {
	return "Set a selected user."
}

func (*SetCmd) Usage() string {
	return ""
}

func (c *SetCmd) SetFlags(f *flag.FlagSet) {
}

func (c *SetCmd) Execute(ctx context.Context, f *flag.FlagSet, args ...interface{}) subcommands.ExitStatus {
	// configからユーザーを選択する
	// 選択したユーザーをgit configに設定する
	// ユーザは引数から取得

	if f.NArg() < 1 {
		log.Printf("You must provide a name for the new directory and an email.")
		return subcommands.ExitFailure
	}
	name := f.Arg(0)

	keys_dir := config.GetKeysDir()

	user_dir := filepath.Join(keys_dir, name)

	if !config.IsDir(user_dir) {
		log.Printf("User %s does not exist.", name)
		return subcommands.ExitFailure
	}

	// json parse
	gitconfig := Config{}
	config_path := filepath.Join(user_dir, "config.json")

	fileData, err := ioutil.ReadFile(config_path)
	if err != nil {
		log.Printf("%v", err)
		return subcommands.ExitFailure
	}
	err = json.Unmarshal(fileData, &gitconfig)
	if err != nil {
		log.Printf("%v", err)
		return subcommands.ExitFailure
	}

	err = config.SetGitUser(gitconfig.Name)
	if err != nil {
		log.Printf("%v", err)
		return subcommands.ExitFailure
	}

	err = config.SetGitEmail(gitconfig.Email)
	if err != nil {
		log.Printf("%v", err)
		return subcommands.ExitFailure
	}

	// keys_dir/user_nameにあるid_ed25519を~/.ssh/id_ed25519にエイリアスする
	// すでに~/.ssh/id_ed25519がある場合は削除する
	// すでに~/.ssh/id_ed25519.pubがある場合は削除する

	// .ssh dirがあるかどうか
	if !config.IsDir(config.GetSSHDir()) {
		// .ssh dirを作成する
		err = os.MkdirAll(config.GetSSHDir(), 0755)
		if err != nil {
			log.Printf("%v", err)
			return subcommands.ExitFailure
		}
	}

	// ~/.ssh/id_ed25519があるかどうか
	ssh_dir := config.GetSSHDir()
	ssh_key_path := filepath.Join(ssh_dir, "id_ed25519")
	ssh_key_pub_path := filepath.Join(ssh_dir, "id_ed25519.pub")

	if config.IsFile(ssh_key_path) {
		// 削除する
		err = config.RemoveFile(ssh_key_path)
		if err != nil {
			log.Printf("%v", err)
			return subcommands.ExitFailure
		}
	}

	if config.IsFile(ssh_key_pub_path) {
		// 削除する
		err = config.RemoveFile(ssh_key_pub_path)
		if err != nil {
			log.Printf("%v", err)
			return subcommands.ExitFailure
		}
	}

	// エイリアスする
	err = config.CreateSymlink(filepath.Join(user_dir, "id_ed25519"), ssh_key_path)
	if err != nil {
		log.Printf("%v", err)
		return subcommands.ExitFailure
	}

	err = config.CreateSymlink(filepath.Join(user_dir, "id_ed25519.pub"), ssh_key_pub_path)
	if err != nil {
		log.Printf("%v", err)
		return subcommands.ExitFailure
	}

	return subcommands.ExitSuccess
}
