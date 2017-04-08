package main

import (
	"math/big"

	"github.com/logic/go-keepassrpc/keepassrpc"
	"github.com/satori/go.uuid"
)

func initSRP(config *configuration) (client *keepassrpc.Client, err error) {
	var value *big.Int

	if config.Username == "" {
		if value, err = keepassrpc.GenKey(32); err != nil {
			return nil, err
		}
		config.Username = uuid.NewV4().String()
	}

	client, err = keepassrpc.NewClient(config.Username, value, config.sessionKey, config.prompt)
	if err != nil {
		return nil, err
	}
	config.sessionKey = client.SessionKey

	if err = config.Save(); err != nil {
		return nil, err
	}

	return client, nil
}
