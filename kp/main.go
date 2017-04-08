package main

import (
	"log"
	"os"

	"github.com/logic/go-keepassrpc/keepassrpc"
)

var config *configuration
var client *keepassrpc.Client

func main() {
	var err error
	config, err = loadConfig()
	if err != nil {
		log.Fatal("loadConfig: ", err)
	}

	client, err = initSRP(config)
	if err != nil {
		log.Fatal("initSRP: ", err)
	}
	defer client.Close()

	ParseCommand(os.Args)
}
