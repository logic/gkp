package keepassrpc

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha1"
	"fmt"
	"io"
	"log"
	"math/big"
	"net/rpc"
	"net/rpc/jsonrpc"

	"github.com/gorilla/websocket"
)

// DebugJSONRPC controls whether unencrypted JSONRPC debugging will be logged
var DebugJSONRPC = false

func pad(data []byte) []byte {
	padding := aes.BlockSize - (len(data) % aes.BlockSize)
	if padding == 0 {
		padding = aes.BlockSize
	}
	return append(data, bytes.Repeat([]byte{byte(padding)}, padding)...)
}

func unpad(data []byte) (output []byte, err error) {
	var dataLen = len(data)
	if dataLen%aes.BlockSize != 0 {
		return output, fmt.Errorf("data's length isn't a multiple of blockSize")
	}
	var paddingBytes = int(data[dataLen-1])
	if paddingBytes > aes.BlockSize || paddingBytes <= 0 {
		return output, fmt.Errorf("invalid padding found: %v", paddingBytes)
	}
	var pad = data[dataLen-paddingBytes : dataLen-1]
	for _, v := range pad {
		if int(v) != paddingBytes {
			return output, fmt.Errorf("invalid padding found")
		}
	}
	output = data[0 : dataLen-paddingBytes]
	return output, nil
}

func hmac(sessionKey, ciphertext, iv []byte) []byte {
	skh := sha1.Sum(sessionKey)
	mac := sha1.New()
	mac.Write(skh[:])
	mac.Write(ciphertext)
	mac.Write(iv)
	return mac.Sum(nil)
}

func encrypt(sessionKey *big.Int, msg []byte) (*MsgJSONRPC, error) {
	plaintext := pad(msg)

	block, err := aes.NewCipher(sessionKey.Bytes())
	if err != nil {
		return nil, err
	}

	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	ciphertext := make([]byte, len(plaintext))
	mode.CryptBlocks(ciphertext, plaintext)

	mac := hmac(sessionKey.Bytes(), ciphertext, iv)

	return &MsgJSONRPC{
		Message: ciphertext,
		IV:      iv,
		HMAC:    mac,
	}, nil
}

func decrypt(sessionKey *big.Int, msg *MsgJSONRPC) ([]byte, error) {
	mac := hmac(sessionKey.Bytes(), msg.Message, msg.IV)
	if bytes.Compare(mac, msg.HMAC) != 0 {
		return nil, fmt.Errorf("HMAC authentication failed")
	}

	block, err := aes.NewCipher(sessionKey.Bytes())
	if err != nil {
		return nil, err
	}

	plaintext := make([]byte, len(msg.Message))
	mode := cipher.NewCBCDecrypter(block, msg.IV)
	mode.CryptBlocks(plaintext, msg.Message)

	return unpad(plaintext)
}

// JSONRPCHandle is our io.ReadWriteCloser implementaion for KeePassRPC crypto
type JSONRPCHandle struct {
	sessionKey *big.Int
	ws         *websocket.Conn
	outbuf     []byte
}

func (ctx *JSONRPCHandle) Write(buf []byte) (int, error) {
	// TODO: this doesn't properly handle messages split over multiple
	// calls to Write. In the real world, this doesn't matter, but it
	// feels like a matter of correctness.
	if DebugJSONRPC {
		log.Print(">>> [JSON-RPC] ", string(buf))
	}

	crypted, err := encrypt(ctx.sessionKey, buf)
	if err != nil {
		return 0, err
	}
	msg := &Message{
		Protocol: "jsonrpc",
		Version:  ProtocolVersion(),
		JSONRPC:  crypted,
	}
	if err := WriteMessage(ctx.ws, msg); err != nil {
		return 0, err
	}
	return len(buf), nil
}

func (ctx *JSONRPCHandle) Read(buf []byte) (int, error) {
	if ctx.outbuf != nil {
		return ctx.popBytes(buf)
	}

	msg, err := ReadMessage(ctx.ws)
	if err != nil {
		return 0, err
	}

	ctx.outbuf, err = decrypt(ctx.sessionKey, msg.JSONRPC)
	if err != nil {
		return 0, err
	}

	if DebugJSONRPC {
		log.Print("<<< [JSON-RPC] ", string(ctx.outbuf))
	}

	return ctx.popBytes(buf)
}

func (ctx *JSONRPCHandle) popBytes(buf []byte) (int, error) {
	var copied int
	if len(ctx.outbuf) > len(buf) {
		copy(buf, ctx.outbuf[:len(buf)])
		copied = len(buf)
		ctx.outbuf = ctx.outbuf[len(buf):]
	} else {
		copy(buf, ctx.outbuf)
		copied = len(ctx.outbuf)
		ctx.outbuf = nil
	}
	return copied, nil
}

// Close closes
func (ctx *JSONRPCHandle) Close() error {
	return fmt.Errorf("Unimplemented")
}

// JSONRPCContext is a wrapper for our websocket that encrypts and decrypts
// automatically
type JSONRPCContext struct {
	c *Client
	r *rpc.Client
}

// DispatchJSONRPC handles setup protocol packets from the server
func DispatchJSONRPC(c *Client, jsonrpc *MsgJSONRPC) error {
	return fmt.Errorf("jsonrpc unimplemented")
}

// EstablishJSONRPCSession sets up our JSON-RPC session
func EstablishJSONRPCSession(c *Client) {
	if c.JSONRPCCtx == nil {
		h := &JSONRPCHandle{
			sessionKey: c.SessionKey,
			ws:         c.WS,
		}
		c.JSONRPCCtx = &JSONRPCContext{
			c: c,
			r: jsonrpc.NewClient(h),
		}
	}
}
