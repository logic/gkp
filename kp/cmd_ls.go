package main

import (
	"flag"
	"fmt"
	"sort"

	"github.com/logic/go-keepassrpc/keepassrpc"
)

type cmdList struct {
	fs      *flag.FlagSet
	long    bool
	recurse bool
}

func (cmd *cmdList) FlagSet() *flag.FlagSet {
	return cmd.fs
}

func (cmd *cmdList) Help() string {
	return "List KeePass entries"
}

type byTitleGroup []keepassrpc.Group

func (g byTitleGroup) Len() int           { return len(g) }
func (g byTitleGroup) Swap(i, j int)      { g[i], g[j] = g[j], g[i] }
func (g byTitleGroup) Less(i, j int) bool { return g[i].Title < g[j].Title }

type byTitleEntry []keepassrpc.Entry

func (e byTitleEntry) Len() int           { return len(e) }
func (e byTitleEntry) Swap(i, j int)      { e[i], e[j] = e[j], e[i] }
func (e byTitleEntry) Less(i, j int) bool { return e[i].Title < e[j].Title }

func cmdListPrintSingleEntry(e *keepassrpc.Entry, long bool) {
	fmt.Printf("%s %s\n", e.UniqueID, e.Title)
}

func cmdListPrintSingleGroup(g *keepassrpc.Group, long bool) {
	fmt.Printf("%s %s/\n", g.UniqueID, g.Title)
}

func cmdListPrintGroup(g *keepassrpc.Group, prefix string, recurse, long bool) error {
	cg, err := client.GetChildGroups(g.UniqueID)
	if err != nil {
		return err
	}
	sort.Sort(byTitleGroup(cg))

	ce, err := client.GetChildEntries(g.UniqueID)
	if err != nil {
		return err
	}
	sort.Sort(byTitleEntry(ce))

	if long {
		for _, c := range cg {
			cmdListPrintSingleGroup(&c, long)
		}
		for _, c := range ce {
			cmdListPrintSingleEntry(&c, long)
		}
	} else {
		names := make([]string, len(cg)+len(ce))
		for i, c := range cg {
			names[i] = fmt.Sprintf("%s/", c.Title)
		}
		for i, c := range ce {
			names[len(cg)+i] = c.Title
		}
		sort.Strings(names)
		for _, n := range names {
			fmt.Println(n)
		}
	}

	if recurse {
		for i := range cg {
			c := cg[i]
			newprefix := fmt.Sprintf("%s/%s", prefix, c.Title)
			fmt.Printf("\n%s:\n", newprefix)
			cmdListPrintGroup(&c, newprefix, recurse, long)
		}
	}
	return nil
}

func (cmd *cmdList) Run(args []string) (err error) {
	var g *keepassrpc.Group
	if len(args) == 0 {
		if g, err = client.GetRoot(); err != nil {
			return err
		}
	} else {
		return fmt.Errorf("specifying custom root unimplemented")
	}

	return cmdListPrintGroup(g, g.Title, cmd.recurse, cmd.long)
}

func init() {
	cmd := &cmdList{
		fs: flag.NewFlagSet("ls", flag.ExitOnError),
	}
	cmd.fs.BoolVar(&cmd.long, "l", false,
		"list entries in long form")
	cmd.fs.BoolVar(&cmd.recurse, "R", false,
		"list entries recursively")
	subcommands["ls"] = cmd
}
