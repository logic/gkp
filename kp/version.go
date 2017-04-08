package main

import "fmt"

var version = "development version"
var timestamp = "unknown"

func versionString() string {
	return fmt.Sprintf("kp %s (build date %s)", version, timestamp)
}
