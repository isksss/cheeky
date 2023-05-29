package main

import (
	"context"
	"flag"
	"os"

	"github.com/google/subcommands"
	"github.com/isksss/cheeky/commands"
)

func main() {
	run()
}

func run() {
	//todo: subcommands
	ctx := context.Background()

	os.Exit(int(subcommands.Execute(ctx)))
}

func init() {
	subcommands.Register(subcommands.CommandsCommand(), "help")
	subcommands.Register(subcommands.FlagsCommand(), "help")
	subcommands.Register(subcommands.HelpCommand(), "help")

	subcommands.Register(&commands.InitCmd{}, "")

	flag.Parse()
}
