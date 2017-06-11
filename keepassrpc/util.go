package keepassrpc

import (
	"crypto/rand"
	"io"
	"math/big"
)

// Our supported KeePassRPC protocol version
var protocolVersion = []uint8{1, 7, 0}

// ProtocolVersion squashes a ProtocolVersion to an int
func ProtocolVersion() uint32 {
	value := uint32(protocolVersion[0])
	value = (value << 8) | uint32(protocolVersion[1])
	value = (value << 8) | uint32(protocolVersion[2])
	return value
}

// GenKey generates a random 32-byte array
func GenKey(len int) (*big.Int, error) {
	bytes := make([]byte, len)
	_, err := io.ReadFull(rand.Reader, bytes)
	if err != nil {
		return nil, err
	}
	return new(big.Int).SetBytes(bytes), nil
}
