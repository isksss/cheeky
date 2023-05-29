package commands

import (
	"context"
	"flag"

	"github.com/google/subcommands"
)

type InitCmd struct{}

func (*InitCmd) Name() string {
	return "init"
}

func (*InitCmd) Synopsis() string {
	return "Initialize a new cheeky project"
}

func (*InitCmd) Usage() string {
	return ""
}

func (c *InitCmd) SetFlags(f *flag.FlagSet) {
}

func (c *InitCmd) Execute(ctx context.Context, f *flag.FlagSet, args ...interface{}) subcommands.ExitStatus {
	return subcommands.ExitSuccess
}
