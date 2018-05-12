package cli

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/bsm/reason/classification/hoeffding"
	"github.com/bsm/reason/internal/iox"
	"github.com/google/subcommands"
)

// ExportCmd exports trees
type ExportCmd struct {
	format string
	output string
}

func (*ExportCmd) Name() string     { return "cls.hoeffding.export" }
func (*ExportCmd) Synopsis() string { return "Export a tree." }
func (*ExportCmd) Usage() string {
	return `cls.hoeffding.export [-f FORMAT] [-o OUTPUT] INPUT

  Exports a hoeffding tree to a file or stdout.
  Inputs may be plain or gzipped files, or specify '-' to read from stdin.

`
}

func (c *ExportCmd) SetFlags(f *flag.FlagSet) {
	f.StringVar(&c.format, "f", "txt", "The format: txt or dot.")
	f.StringVar(&c.output, "o", "-", "The output, either a file path or '-' for stdout.")
}

func (c *ExportCmd) Execute(ctx context.Context, f *flag.FlagSet, args ...interface{}) subcommands.ExitStatus {
	ifname := f.Arg(0)
	if ifname == "" {
		fmt.Fprintf(os.Stderr, "[ERROR] no input specified\n")
		return subcommands.ExitUsageError
	}

	switch c.format {
	case "txt", "dot":
	default:
		fmt.Fprintf(os.Stderr, "[ERROR] unknown format %q\n", c.format)
		return subcommands.ExitUsageError
	}

	input, err := iox.Open(ifname)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[ERROR] unable to open input: %s\n", err.Error())
		return subcommands.ExitFailure
	}
	defer input.Close()

	tree, err := hoeffding.Load(input, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[ERROR] unable to load tree: %s\n", err.Error())
		return subcommands.ExitFailure
	}

	output, err := iox.Create(c.output)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[ERROR] unable to create output: %s\n", err.Error())
		return subcommands.ExitFailure
	}
	defer output.Close()

	switch c.format {
	case "txt":
		_, err = tree.WriteText(output)
	case "dot":
		_, err = tree.WriteDOT(output)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "[ERROR] unable to export tree: %s\n", err.Error())
		return subcommands.ExitFailure
	}

	if err := output.Close(); err != nil {
		fmt.Fprintf(os.Stderr, "[ERROR] unable to flush output: %s\n", err.Error())
		return subcommands.ExitFailure
	}

	return subcommands.ExitSuccess
}
