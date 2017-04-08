package main

import (
	"flag"
	"fmt"
	"os"
)

type command interface {
	FlagSet() *flag.FlagSet
	Run([]string) error
	Help() string
}

var subcommands = map[string]command{}

func globalHelp() {
	flag.Usage()
	fmt.Println("Valid subcommands:")
	for name, cmd := range subcommands {
		fmt.Printf("    %-10s %s\n", name, cmd.Help())
	}
	os.Exit(1)
}

// ParseCommand takes a command line and works out what to do next
func ParseCommand(args []string) {
	if len(args) < 2 {
		globalHelp()
	}
	if fs, ok := subcommands[args[1]]; ok {
		fs.FlagSet().Parse(args[2:])
		if err := fs.Run(fs.FlagSet().Args()); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	} else {
		globalHelp()
	}
}
