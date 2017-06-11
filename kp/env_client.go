package main

import "github.com/logic/gkp/keepassrpc"

type envDebugClient struct{}

func (env *envDebugClient) Trigger(value string) error {
	keepassrpc.DebugClient = true
	return nil
}

func (env *envDebugClient) Help() string {
	return "Debug encryption-establishing client protocol"
}

func init() {
	envvars["KEEPASSRPC_DEBUG_CLIENT"] = &envDebugClient{}
}
