package main

import (
	"flag"
	"fmt"
)

type cmdVersion struct {
	fs *flag.FlagSet
}

func (cmd *cmdVersion) FlagSet() *flag.FlagSet {
	return cmd.fs
}

func (cmd *cmdVersion) Run(args []string) error {
	fmt.Println(versionString())
	return nil
}

func (cmd *cmdVersion) Help() string {
	return "Version queries"
}

func init() {
	cmd := &cmdVersion{
		fs: flag.NewFlagSet("version", flag.ExitOnError),
	}
	subcommands["version"] = cmd
}
