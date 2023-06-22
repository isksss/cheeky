package main

import (
	"context"
	"flag"
	"log"
	"os"

	"github.com/google/subcommands"
	"github.com/isksss/cheeky/commands"
	"github.com/isksss/cheeky/config"
)

func main() {
	run()
}

func run() {
	//todo: subcommands
	parse()

	if !config.IsGitInstalled() {
		log.Fatalln("Git is not installed.")
		os.Exit(1)
	}

	ctx := context.Background()
	os.Exit(int(subcommands.Execute(ctx)))
}

func init() {
	config.MakeDir()
}

func parse() {
	subcommands.Register(subcommands.CommandsCommand(), "help")
	subcommands.Register(subcommands.FlagsCommand(), "help")
	subcommands.Register(subcommands.HelpCommand(), "help")

	subcommands.Register(&commands.InitCmd{}, "")
	subcommands.Register(&commands.NewCmd{}, "")
	subcommands.Register(&commands.SetCmd{}, "")

	flag.Parse()
}
