package keepassrpc

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/gorilla/websocket"
)

// MsgError represents an error in the KeePassRPC protocol
type MsgError struct {
	Code          string   `json:"code"`
	MessageParams []string `json:"messageParams"`
}

// MsgJSONRPC represents various stages of the JSON RPC protocol
type MsgJSONRPC struct {
	Message []byte `json:"message"`
	IV      []byte `json:"iv"`
	HMAC    []byte `json:"hmac"`
}

// Message represents a single complete KeePassRPC message
type Message struct {
	// Mandatory values
	Protocol string `json:"protocol"`
	Version  uint32 `json:"version"`

	JSONRPC    *MsgJSONRPC `json:"jsonrpc,omitempty"`
	Error      *MsgError   `json:"error,omitempty"`
	SRP        *MsgSRP     `json:"srp,omitempty"`
	Key        *MsgKey     `json:"key,omitempty"`
	ClientID   string      `json:"clientTypeID,omitempty"`
	ClientName string      `json:"clientDisplayName,omitempty"`
	ClientDesc string      `json:"clientDisplayDescription,omitempty"`
	Features   []string    `json:"features,omitempty"`
}

// ReadMessage reads a message from the server and JSON-decodes it
func ReadMessage(h *websocket.Conn) (*Message, error) {
	var data Message
	if DebugClient {
		_, msg, err := h.ReadMessage()
		if err != nil {
			return nil, err
		}
		log.Println("<<<", string(msg))
		if err := json.Unmarshal(msg, &data); err != nil {
			return nil, err
		}
	} else {
		if err := h.ReadJSON(&data); err != nil {
			return nil, err
		}
	}

	return &data, nil
}

// WriteMessage JSON-encodes a KeePassRPC message and sends it to the server
func WriteMessage(h *websocket.Conn, msg *Message) error {
	if DebugClient {
		out, err := json.Marshal(msg)
		if err != nil {
			return err
		}
		log.Println(">>>", string(out))
		if err := h.WriteMessage(websocket.TextMessage, out); err != nil {
			return err
		}
	} else {
		if err := h.WriteJSON(&msg); err != nil {
			return err
		}
	}
	return nil
}

// DispatchError handles error protocol packets from the server
func DispatchError(err *MsgError) error {
	return fmt.Errorf("%s: %s", err.Code, strings.Join(err.MessageParams, "\n"))
}
