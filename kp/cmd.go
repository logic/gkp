package main

import (
	"flag"
	"fmt"
	"os"
	"sort"
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

	names := make([]string, len(subcommands))
	i := 0
	for k := range subcommands {
		names[i] = k
		i++
	}
	sort.Strings(names)

	for _, name := range names {
		fmt.Printf("    %-10s %s\n", name, subcommands[name].Help())
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
