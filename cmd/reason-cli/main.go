package main

import (
	"context"
	"flag"
	"os"

	cls_hoeffding "github.com/bsm/reason/classification/hoeffding/cli"
	"github.com/google/subcommands"
)

func main() {
	subcommands.Register(subcommands.HelpCommand(), "")
	subcommands.Register(subcommands.CommandsCommand(), "")
	subcommands.Register(&cls_hoeffding.ExportCmd{}, "")

	flag.Parse()
	ctx := context.Background()
	os.Exit(int(subcommands.Execute(ctx)))
}
