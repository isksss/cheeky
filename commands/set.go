package commands

import (
	"context"
	"flag"

	"github.com/google/subcommands"
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
	return subcommands.ExitSuccess
}
