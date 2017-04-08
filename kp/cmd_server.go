package main

import (
	"flag"
	"fmt"
)

type cmdServer struct {
	fs *flag.FlagSet
}

func (cmd *cmdServer) FlagSet() *flag.FlagSet {
	return cmd.fs
}

func (cmd *cmdServer) Run(args []string) error {
	info, err := client.GetApplicationMetadata()
	if err != nil {
		return err
	}
	fmt.Println("KeePass", info.KeePassVersion)
	if info.IsMono {
		fmt.Println("Mono", info.MonoVersion)
	} else {
		fmt.Println(".NET", info.NETversion)
	}
	fmt.Println(".NET CLR", info.NETCLR)

	about, err := client.SystemAbout()
	if err != nil {
		return err
	}
	fmt.Print(about)

	return nil
}

func (cmd *cmdServer) Help() string {
	return "Information about the running KeePass instance"
}

func init() {
	cmd := &cmdServer{
		fs: flag.NewFlagSet("server", flag.ExitOnError),
	}
	subcommands["server"] = cmd
}
