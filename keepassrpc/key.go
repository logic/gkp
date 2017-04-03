package keepassrpc

import (
	"crypto/sha256"
	"fmt"
)

// MsgKey represents various stages of the challenge/response protocol
type MsgKey struct {
	SecurityLevel int `json:"securityLevel"`

	SC       string `json:"sc,omitempty"` // Server challenge
	CC       string `json:"cc,omitempty"` // Client challenge
	SR       string `json:"sr,omitempty"` // Server response
	CR       string `json:"cr,omitempty"` // Client response
	Username string `json:"username,omitempty"`
}

// KeyContext is key-negotiation-specific context that we track per-client
type KeyContext struct {
	sc string
	cc string
}

// DispatchKey handles challenge/response setup negotiation with the server
func DispatchKey(c *Client, key *MsgKey) error {
	if key.SC != "" {
		return ServerChallenge(c, key)
	}
	if key.SR != "" {
		return ServerResponse(c, key)
	}
	return fmt.Errorf("neither server challenge nor response provided")
}

// EstablishKeySession establishes a session with the KeePassRPC service using
// a previously-negotiated key.
func EstablishKeySession(c *Client) error {
	challenge, err := GenKey(32)
	if err != nil {
		return err
	}
	c.KeyCtx = &KeyContext{cc: challenge.Text(16)}
	defer func() { c.KeyCtx = nil }()

	msg := &Message{
		Protocol:   "setup",
		ClientID:   ClientID,
		ClientName: ClientName,
		ClientDesc: ClientDesc,
		Version:    ProtocolVersion(),
		Key: &MsgKey{
			Username:      c.Username,
			SecurityLevel: 2,
		},
	}

	if err := WriteMessage(c.WS, msg); err != nil {
		return err
	}

	return c.DispatchResponse()
}

// ServerChallenge responds to a challenge issued by the server, and issues
// a challenge of our own.
func ServerChallenge(c *Client, key *MsgKey) error {
	c.KeyCtx.sc = key.SC

	h := sha256.New()
	fmt.Fprintf(h, "1%x%s%s", c.SessionKey, c.KeyCtx.sc, c.KeyCtx.cc)
	response := fmt.Sprintf("%x", h.Sum(nil))

	resp := &Message{
		Protocol: "setup",
		Version:  ProtocolVersion(),
		Key: &MsgKey{
			CC:            c.KeyCtx.cc,
			CR:            response,
			SecurityLevel: 2,
		},
	}
	if err := WriteMessage(c.WS, resp); err != nil {
		return err
	}
	return c.DispatchResponse()
}

// ServerResponse validates the server's response to our challenge.
func ServerResponse(c *Client, key *MsgKey) error {
	h := sha256.New()
	fmt.Fprintf(h, "0%x%s%s", c.SessionKey, c.KeyCtx.sc, c.KeyCtx.cc)
	sr := fmt.Sprintf("%x", h.Sum(nil))

	c.KeyCtx = nil

	if sr != key.SR {
		return fmt.Errorf("Server key does not match")
	}

	return nil
}
