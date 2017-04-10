package keepassrpc

import (
	"fmt"
	"log"
	"math/big"

	"github.com/gorilla/websocket"
)

// DebugClient controls whether protocol debugging will be logged
var DebugClient = false

// DefaultURL is the canonical local location of the KeePassRPC service
const DefaultURL = "ws://127.0.0.1:12546/"

// ClientID is a short ID for our client
const ClientID = "gkp"

// ClientName is a short human-readable client description
const ClientName = "gkp-KeePassRPC"

// ClientDesc is a longer human-readable client description
const ClientDesc = "A Go KeePassRPC implementation"

// Passworder is expected to return either a string password (the nonce provided
// by KeePass) or an error
type Passworder func() (string, error)

// Client represents an operating KeePassRPC session.
type Client struct {
	Username   string
	SessionKey *big.Int
	Value      *big.Int
	Password   Passworder

	WS *websocket.Conn

	SRPCtx     *SRPContext
	KeyCtx     *KeyContext
	JSONRPCCtx *JSONRPCContext
}

// NewClient instantiates a new KeePassRPC client for the given user
func NewClient(username string, value, sessionKey *big.Int, pwd Passworder) (*Client, error) {
	wsc, _, err := websocket.DefaultDialer.Dial(DefaultURL, nil)
	if err != nil {
		return nil, err
	}

	c := &Client{
		Username:   username,
		SessionKey: sessionKey,
		Value:      value,
		Password:   pwd,
		WS:         wsc,
	}

	if err := c.EstablishSession(); err != nil {
		return nil, err
	}

	return c, nil
}

// DispatchResponse is the general server response handler
func (c *Client) DispatchResponse() error {
	msg, err := ReadMessage(c.WS)
	if err != nil {
		return err
	}

	// Ugh.
	if msg.Error != nil {
		return DispatchError(msg.Error)
	}

	switch msg.Protocol {
	case "setup":
		if msg.SRP != nil && msg.Key == nil {
			return DispatchSRP(c, msg.SRP)
		}
		if msg.Key != nil && msg.SRP == nil {
			return DispatchKey(c, msg.Key)
		}
		return fmt.Errorf("invalid setup")
	case "jsonrpc":
		return DispatchJSONRPC(c, msg.JSONRPC)
	case "error":
		return DispatchError(msg.Error)
	}
	return fmt.Errorf("Unknown protocol '%s'", msg.Protocol)
}

// EstablishSession starts a new KeePassRPC session, either via a new SRP
// negotiation, or via challenge/response with an established key.
func (c *Client) EstablishSession() error {
	if c.SessionKey != nil {
		if err := EstablishKeySession(c); err != nil {
			// Treat an initial failure as temporary; the key might
			// have simply expired or been revoked, so we need to
			// go through a new SRP phase.
			if DebugClient {
				log.Println("Couldn't establish session with existing key:", err)
			}
			c.SessionKey = nil
		}
	}

	// If we don't have a valid session key (or it was just rejected), try
	// a fresh SRP session to negotiate a new session key.
	if c.SessionKey == nil {
		if err := EstablishSRPSession(c); err != nil {
			return err
		}
	}

	// Post-authentication, establish our JSON-RPC session.
	EstablishJSONRPCSession(c)
	return nil
}

// Close closes the client's underlying websocket, if it exists
func (c *Client) Close() {
	if c.WS != nil {
		c.WS.WriteMessage(websocket.CloseMessage,
			websocket.FormatCloseMessage(
				websocket.CloseNormalClosure, "goodbye"))
		c.WS.Close()
	}
}
