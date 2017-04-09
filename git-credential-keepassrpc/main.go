package main

import (
	"bufio"
	"fmt"
	"github.com/logic/go-keepassrpc/keepassrpc/cli"
	"io"
	"net/url"
	"os"
	"strings"
)

// ReadCredential reads a git-credential formatted input block into a URL
func ReadCredential(f io.Reader) *url.URL {
	creds := map[string]string{}
	fscanner := bufio.NewScanner(f)
	for fscanner.Scan() {
		t := strings.SplitN(fscanner.Text(), "=", 2)
		creds[strings.TrimSpace(t[0])] = strings.TrimSpace(t[1])
	}

	var u *url.URL
	if uu, ok := creds["url"]; ok {
		u, _ = url.Parse(uu)
	}
	if u == nil {
		u = &url.URL{}
		if p, ok := creds["protocol"]; ok {
			u.Scheme = p
		}
		if h, ok := creds["host"]; ok {
			u.Host = h
		}
		if h, ok := creds["host"]; ok {
			u.Host = h
		}
		if path, ok := creds["path"]; ok {
			u.Path = path
		}
		if username, ok := creds["username"]; ok {
			if password, ok := creds["password"]; ok {
				u.User = url.UserPassword(username, password)
			} else {
				u.User = url.User(username)
			}
		}
	}

	return u
}

// GetCredentials retrieves a credential based on supplied data.
func GetCredentials(u *url.URL) {
	config, err := cli.LoadConfig()
	if err != nil {
		return
	}

	// TODO: is there a reasonable way to prompt the user here?
	client, err := cli.Dial(config, nil)
	if err != nil {
		return
	}

	s := client.NewSearch()
	s.AddURL(u.String())
	entries, err := s.Execute()
	if err != nil {
		return
	}

	if len(entries) > 0 {
		e := entries[0]
		u.User = url.UserPassword(e.Username(), e.Password())
	}
}

// StoreCredentials stores an update to the supplied credentials.
func StoreCredentials(u *url.URL) {
	// do nothing right now
}

// EraseCredentials erases the referenced credentials.
func EraseCredentials(u *url.URL) {
	// do nothing right now
}

func main() {
	if len(os.Args) != 2 {
		panic("Need a single operation (get/store/erase) as argument")
	}

	u := ReadCredential(os.Stdin)

	switch os.Args[1] {
	case "get":
		GetCredentials(u)
	case "store":
		StoreCredentials(u)
	case "erase":
		EraseCredentials(u)
	default:
		panic(fmt.Sprintf("Unknown operation '%s'", os.Args[1]))
	}

	fmt.Printf("url=%s\n", u)
}
