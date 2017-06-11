package main

import "os"

type envvar interface {
	Trigger(string) error
	Help() string
}

var envvars = map[string]envvar{}

// ParseEnvironment walks the list of registered environment variables and
// calls their trigger functions.
func ParseEnvironment() {
	for name, action := range envvars {
		if value, ok := os.LookupEnv(name); ok {
			action.Trigger(value)
		}
	}
}
