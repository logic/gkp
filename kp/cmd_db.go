package main

import (
	"flag"
	"fmt"
)

type cmdDatabase struct {
	fs           *flag.FlagSet
	CloseCurrent bool
}

func (cmd *cmdDatabase) FlagSet() *flag.FlagSet {
	return cmd.fs
}

func (cmd *cmdDatabase) Run(args []string) error {
	if len(args) == 1 {
		return client.ChangeDatabase(args[0], cmd.CloseCurrent)
	}

	if len(args) > 0 {
		return fmt.Errorf("You can only specify a single database")
	}

	config, err := client.GetCurrentKFConfig()
	if err != nil {
		return err
	}
	db, err := client.GetDatabaseName()
	if err != nil {
		return err
	}
	dbFname, err := client.GetDatabaseFileName()
	if err != nil {
		return err
	}

	fmt.Println("Databases: (*** = active)")
	for _, d := range config.KnownDatabases {
		if d == dbFname {
			fmt.Printf("*** %s (named '%s')\n", d, db)
		} else {
			fmt.Println("   ", d)
		}
	}
	fmt.Println("Autocommit on change?", config.AutoCommit)

	return nil
}

func (cmd *cmdDatabase) Help() string {
	return "Information about the running KeePass instance"
}

func init() {
	cmd := &cmdDatabase{
		fs: flag.NewFlagSet("db", flag.ExitOnError),
	}
	cmd.fs.BoolVar(&cmd.CloseCurrent, "close", true,
		"Close the current database when switching to a new one")
	subcommands["db"] = cmd
}
