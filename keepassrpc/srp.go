package keepassrpc

import (
	"crypto/sha1"
	"crypto/sha256"
	"errors"
	"fmt"
	"math/big"
)

// Generator is the SRP generator that KeePassRPC uses
var Generator = big.NewInt(2)

// Prime is the large safe prime that KeePassRPC uses for SRP negotiation
var Prime = new(big.Int).SetBytes([]byte("" +
	"\xd4\xc7\xf8\xa2\xb3\x2c\x11\xb8\xfb\xa9\x58\x1e\xc4\xba\x4f\x1b" +
	"\x04\x21\x56\x42\xef\x73\x55\xe3\x7c\x0f\xc0\x44\x3e\xf7\x56\xea" +
	"\x2c\x6b\x8e\xeb\x75\x5a\x1c\x72\x30\x27\x66\x3c\xaa\x26\x5e\xf7" +
	"\x85\xb8\xff\x6a\x9b\x35\x22\x7a\x52\xd8\x66\x33\xdb\xdf\xca\x43" +
	""))

// MsgSRP represents various stages of the SRP protocol
type MsgSRP struct {
	// Mandatory values
	Stage         string `json:"stage"`
	SecurityLevel int    `json:"securityLevel"`

	A  string `json:",omitempty"`  // Client public value
	B  string `json:",omitempty"`  // Server public value
	I  string `json:",omitempty"`  // Identity (username)
	M  string `json:",omitempty"`  // Evidence
	M2 string `json:",omitempty"`  // Server Evidence
	S  string `json:"s,omitempty"` // Salt
}

// SRPContext represents the client side of our SRP negotiation
type SRPContext struct {
	Public   *big.Int // A
	Private  *big.Int // a
	Server   *big.Int // B
	Password string   // P
	Salt     string   // s

	// Memoized values
	_x  *big.Int
	_u  *big.Int
	_S  *big.Int
	_M  *big.Int
	_M2 *big.Int
}

// Multiplier calculates the SRP multiplier to be used, which is just a SHA1
// hash of our chosen prime (N), plus our generator padded to the length of N.
func (c *SRPContext) Multiplier() *big.Int {
	h := sha1.New()
	h.Write(Prime.Bytes())
	padLen := len(Prime.Bytes()) - len(Generator.Bytes())
	pad := make([]byte, padLen)
	h.Write(pad)
	h.Write(Generator.Bytes())
	return new(big.Int).SetBytes(h.Sum(nil))
}

// SetServer validates and stores the server public value
func (c *SRPContext) SetServer(B *big.Int) error {
	// RFC5054 2.5.3 Server Key Exchange
	// The client MUST abort the handshake with an "illegal_parameter"
	// alert if B % N = 0.
	t := new(big.Int)
	if t.Mod(B, Prime); t.Cmp(big.NewInt(0)) == 0 {
		return errors.New("illegal parameter")
	}

	c.Server = new(big.Int).Set(B)
	return nil
}

func (c *SRPContext) privateKey() *big.Int {
	if c._x == nil {
		h := sha256.New()
		fmt.Fprintf(h, "%s%s", c.Salt, c.Password)
		c._x = new(big.Int).SetBytes(h.Sum(nil))
	}
	return c._x
}

func (c *SRPContext) scramblingParameter() *big.Int {
	if c._u == nil {
		h := sha256.New()
		fmt.Fprintf(h, "%X%X", c.Public, c.Server)
		c._u = new(big.Int).SetBytes(h.Sum(nil))
	}
	return c._u
}

func (c *SRPContext) premasterSecret() *big.Int {
	if c._S == nil {
		x := c.privateKey()

		c._S = new(big.Int)
		c._S.Exp(Generator, x, Prime)
		c._S.Mul(c.Multiplier(), c._S)
		c._S.Sub(c.Server, c._S)

		exp := new(big.Int)
		exp.Mul(c.scramblingParameter(), x)
		exp.Add(c.Private, exp)

		c._S.Exp(c._S, exp, Prime)
	}

	return c._S
}

// Verifier returns the verifier value for the current user
func (c *SRPContext) Verifier() (*big.Int, error) {
	if c.Password == "" || c.Salt == "" {
		return nil, errors.New("premature generation of verifier")
	}
	return new(big.Int).Exp(Generator, c.privateKey(), Prime), nil
}

// Evidence calculates client-side evidence for submission to the server
func (c *SRPContext) Evidence() (*big.Int, error) {
	if c._M == nil {
		if c.Server == nil {
			return nil, errors.New("premature generation of evidence")
		}

		S := c.premasterSecret()

		h := sha256.New()
		fmt.Fprintf(h, "%X%X%X", c.Public, c.Server, S)
		c._M = new(big.Int).SetBytes(h.Sum(nil))
	}

	return c._M, nil
}

// ServerEvidence calculates server-side evidence for comparison
func (c *SRPContext) ServerEvidence() (*big.Int, error) {
	if c._M2 == nil {
		M, err := c.Evidence()
		if err != nil {
			return nil, err
		}

		S := c.premasterSecret()

		h := sha256.New()
		fmt.Fprintf(h, "%X%x%X", c.Public, M, S)
		c._M2 = new(big.Int).SetBytes(h.Sum(nil))
	}

	return c._M2, nil
}

// SessionKey returns a reusable session key for use in future connections
func (c *SRPContext) SessionKey() *big.Int {
	h := sha256.New()
	fmt.Fprintf(h, "%X", c.premasterSecret())
	return new(big.Int).SetBytes(h.Sum(nil))
}

// DispatchSRP handles SRP setup negotiation with the server
func DispatchSRP(c *Client, srp *MsgSRP) error {
	switch srp.Stage {
	case "identifyToClient":
		return IdentifyToClient(c, srp)
	case "proofToClient":
		return ProofToClient(c, srp)
	}
	return fmt.Errorf("unknown SRP stage '%s'", srp.Stage)
}

// EstablishSRPSession establishes a session with the KeePassRPC service using
// a fresh SRP negotiation.
func EstablishSRPSession(c *Client) error {
	ctx := &SRPContext{
		Private:  new(big.Int).Set(c.Value),
		Public:   new(big.Int).Exp(Generator, c.Value, Prime),
		Salt:     "",
		Server:   nil,
		Password: "",
	}

	c.SRPCtx = ctx
	defer func() { c.SRPCtx = nil }()

	msg := &Message{
		Protocol:   "setup",
		ClientID:   ClientID,
		ClientName: ClientName,
		ClientDesc: ClientDesc,
		Version:    ProtocolVersion(),
		SRP: &MsgSRP{
			Stage:         "identifyToServer",
			I:             c.Username,
			A:             fmt.Sprintf("%X", c.SRPCtx.Public),
			SecurityLevel: 2,
		},
	}

	if err := WriteMessage(c.WS, msg); err != nil {
		return err
	}

	if err := c.DispatchResponse(); err != nil {
		return err
	}

	c.SessionKey = ctx.SessionKey()
	return nil
}

// IdentifyToClient is the initial response from the KeePassRPC service
func IdentifyToClient(c *Client, srp *MsgSRP) error {
	password, err := c.Password()
	if err != nil {
		return err
	}
	c.SRPCtx.Password = password

	// Pad value if needed.
	if (len(srp.B) % 2) != 0 {
		srp.B = "0" + srp.B
	}

	B, ok := new(big.Int).SetString(srp.B, 16)
	if !ok {
		return fmt.Errorf("Could not decode server-provided B")
	}
	if err := c.SRPCtx.SetServer(B); err != nil {
		return err
	}

	c.SRPCtx.Salt = srp.S

	return ProofToServer(c)
}

// ProofToServer issues our proof to the KeePassRPC server
func ProofToServer(c *Client) error {
	M, err := c.SRPCtx.Evidence()
	if err != nil {
		return err
	}
	msg := &Message{
		Protocol: "setup",
		Version:  ProtocolVersion(),
		SRP: &MsgSRP{
			Stage:         "proofToServer",
			M:             fmt.Sprintf("%X", M),
			SecurityLevel: 2,
		},
	}
	if err := WriteMessage(c.WS, msg); err != nil {
		return err
	}
	return c.DispatchResponse()
}

// ProofToClient verifies proof provided by the KeePassRPC server
func ProofToClient(c *Client, srp *MsgSRP) error {
	M2, ok := new(big.Int).SetString(srp.M2, 16)
	if !ok {
		return fmt.Errorf("Could not decode server-provided M2")
	}

	ourM2, err := c.SRPCtx.ServerEvidence()
	if err != nil {
		return err
	}

	if M2.Cmp(ourM2) != 0 {
		return fmt.Errorf("Server-provided evidence does not match")
	}

	return nil
}
