package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/logic/go-keepassrpc/keepassrpc"
)

type cmdSearch struct {
	fs       *flag.FlagSet
	ShowAll  bool
	UniqueID string
	URLs     bool
	Results  int
}

func (cmd *cmdSearch) FlagSet() *flag.FlagSet {
	return cmd.fs
}

func (cmd *cmdSearch) Help() string {
	return "Search KeePass for free-text search terms, unique IDs, and URLs"
}

func (cmd *cmdSearch) Run(args []string) error {
	if cmd.UniqueID != "" && len(args) != 0 {
		return fmt.Errorf("must specify a single unique ID")
	}

	s := client.NewSearch()
	if cmd.UniqueID != "" {
		s.UniqueID = cmd.UniqueID
	} else if cmd.URLs {
		s.UnsanitizedURLs = append(s.UnsanitizedURLs, args...)
	} else {
		s.FreeTextSearch = strings.Join(args, " ")
	}
	entries, err := s.Execute()
	if err != nil {
		return err
	}
	if len(entries) == 0 {
		fmt.Println("No entries found.")
		return nil
	}

	countdown := cmd.Results
	if countdown == 0 {
		countdown = len(entries)
	}
	for _, e := range entries {
		if countdown < 1 {
			break
		}
		fmt.Println()
		fmt.Print(e.Title)
		switch e.MatchAccuracy {
		case keepassrpc.MatchAccuracyBest:
			fmt.Println(" [best match] ")
		case keepassrpc.MatchAccuracyClose:
			fmt.Println(" [close match]")
		case keepassrpc.MatchAccuracyHostnameAndPort:
			fmt.Println(" [hostname and port match]")
		case keepassrpc.MatchAccuracyHostname:
			fmt.Println(" [hostname match]")
		case keepassrpc.MatchAccuracyDomain:
			fmt.Println(" [domain match]")
		case keepassrpc.MatchAccuracyNone:
			fmt.Println()
		default:
			fmt.Printf(" [unknown match result: %d]\n", e.MatchAccuracy)
		}
		fmt.Println("UUID:", e.UniqueID)
		fmt.Println("URLs:")
		for _, u := range e.URLs {
			fmt.Println("   ", u)
		}
		fmt.Println("Form fields:")
		for _, f := range e.FormFieldList {
			switch f.Type {
			case keepassrpc.FFTradio:
				fmt.Print("🔘")
			case keepassrpc.FFTusername:
				fmt.Print("👩")
			case keepassrpc.FFTpassword:
				fmt.Print("🔒")
			case keepassrpc.FFTselect:
				fmt.Print("▼")
			case keepassrpc.FFTcheckbox:
				fmt.Print("✓")
			case keepassrpc.FFTtext:
				fmt.Print("🗎")
			default:
				fmt.Print("�")
			}
			if f.DisplayName != "" {
				fmt.Print("\t", f.DisplayName)
				if f.DisplayName != f.Name {
					fmt.Print(" (", f.Name)
					if f.Name != f.ID {
						fmt.Print(", ", f.ID)
					}
					fmt.Print(")")
				}
			} else {
				fmt.Print("\t[no name]")
			}
			if f.Value != "" {
				fmt.Print(": ")
				if f.Type == keepassrpc.FFTcheckbox {
					switch f.Value {
					case "KEEFOX_CHECKED_FLAG_FALSE":
						fmt.Print("☐")
					case "KEEFOX_CHECKED_FLAG_TRUE":
						fmt.Print("☑")
					default:
						fmt.Print(f.Value)
					}
				} else if f.Type != keepassrpc.FFTpassword || cmd.ShowAll {
					fmt.Print(f.Value)
				} else {
					fmt.Print("********")
				}
			}
			fmt.Println()
		}

		countdown--
	}

	return nil
}

func init() {
	cmd := &cmdSearch{
		fs: flag.NewFlagSet("search", flag.ExitOnError),
	}
	cmd.fs.BoolVar(&cmd.ShowAll, "showall", false,
		"Display password fields")
	cmd.fs.StringVar(&cmd.UniqueID, "uuid", "",
		"Search for a single unique UUID")
	cmd.fs.BoolVar(&cmd.URLs, "urls", false,
		"Treat arguments as URLs instead of free-text search terms")
	cmd.fs.IntVar(&cmd.Results, "n", 0,
		"Number of results to return (0 = everything)")
	subcommands["search"] = cmd
}
