package main

import (
	"log"
	"os"

	"github.com/logic/gkp/keepassrpc"
	"github.com/logic/gkp/keepassrpc/cli"
)

var config *cli.Configuration
var client *keepassrpc.Client

func main() {
	var err error
	config, err = cli.LoadConfig()
	if err != nil {
		log.Fatal("loadConfig: ", err)
	}

	client, err = cli.Dial(config, cli.Prompt)
	if err != nil {
		log.Fatal("initSRP: ", err)
	}
	defer client.Close()

	ParseCommand(os.Args)
}
