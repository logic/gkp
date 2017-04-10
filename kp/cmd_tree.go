package main

import (
	"flag"
	"fmt"

	"github.com/logic/gkp/keepassrpc"
)

type cmdTree struct {
	fs *flag.FlagSet
}

func (cmd *cmdTree) FlagSet() *flag.FlagSet {
	return cmd.fs
}

func (cmd *cmdTree) Help() string {
	return "List all available KeePass entries as a visual hierarchy"
}

func cmdTreePrintSingle(name string, prefixes []bool, last bool) {
	for _, p := range prefixes {
		if p {
			fmt.Print("    ")
		} else {
			fmt.Print("│   ")
		}
	}

	thisPrefix := "├──"
	if last {
		thisPrefix = "└──"
	}

	fmt.Println(thisPrefix, name)
}

func cmdTreePrintGroup(g *keepassrpc.Group, prefixes []bool) error {
	cg, err := client.GetChildGroups(g.UniqueID)
	if err != nil {
		return err
	}
	ce, err := client.GetChildEntries(g.UniqueID)
	if err != nil {
		return err
	}

	for i := range cg {
		c := cg[i]
		last := (i >= len(cg)-1) && (len(ce) == 0)
		cmdTreePrintSingle(c.Title, prefixes, last)

		newprefixes := append(prefixes, last)
		cmdTreePrintGroup(&c, newprefixes)
	}
	for i := range ce {
		c := ce[i]
		cmdTreePrintSingle(c.Title, prefixes, (i >= len(ce)-1))
	}
	return nil
}

func cmdTreePrintTree(root *keepassrpc.Group) error {
	fmt.Println(root.Title)
	return cmdTreePrintGroup(root, nil)
}

func (cmd *cmdTree) Run(args []string) (err error) {
	var g *keepassrpc.Group
	if len(args) == 0 {
		if g, err = client.GetRoot(); err != nil {
			return err
		}
	} else {
		return fmt.Errorf("specifying custom root unimplemented")
	}

	fmt.Println(g.Title)
	return cmdTreePrintGroup(g, nil)
}

func init() {
	cmd := &cmdTree{
		fs: flag.NewFlagSet("tree", flag.ExitOnError),
	}
	subcommands["tree"] = cmd
}
