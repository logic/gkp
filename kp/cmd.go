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

	names := make([]string, len(subcommands))
	i := 0
	for k := range subcommands {
		names[i] = k
		i++
	}
	sort.Strings(names)

	fmt.Println("Valid subcommands:")
	for _, name := range names {
		fmt.Printf("    %-10s %s\n", name, subcommands[name].Help())
	}

	if len(envvars) != 0 {
		fmt.Println("\nValid environment variables:")
		for name, action := range envvars {
			fmt.Printf("    %-25s %s\n", name, action.Help())
		}
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
