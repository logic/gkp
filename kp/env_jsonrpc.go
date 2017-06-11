package main

import "github.com/logic/gkp/keepassrpc"

type envDebugJSONRPC struct{}

func (env *envDebugJSONRPC) Trigger(value string) error {
	keepassrpc.DebugJSONRPC = true
	return nil
}

func (env *envDebugJSONRPC) Help() string {
	return "Debug post-encryption JSON-RPC protocol"
}

func init() {
	envvars["KEEPASSRPC_DEBUG_JSONRPC"] = &envDebugJSONRPC{}
}
