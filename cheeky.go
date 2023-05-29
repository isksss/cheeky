package main

import (
	"context"
	"flag"
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

	flag.Parse()
}
