package main

import (
	"github.com/logic/go-keepassrpc/keepassrpc"
	"github.com/satori/go.uuid"
)

func initSRP(config *configuration) (err error) {
	if config.Username == "" {
		if config.Value, err = keepassrpc.GenKey(32); err != nil {
			return err
		}
		config.Username = uuid.NewV4().String()
	}

	if config.client, err = keepassrpc.NewClient(config.Username,
		config.Value, config.SessionKey, config.prompt); err != nil {
		return err
	}
	config.SessionKey = config.client.SessionKey

	if err := config.Save(); err != nil {
		return err
	}

	return nil
}

func init() {
	keepassrpc.DebugJSONRPC = true
}
