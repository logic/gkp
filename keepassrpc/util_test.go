package keepassrpc

import "testing"

func TestVersion(t *testing.T) {
	protocolVersion = []uint8{1, 2, 3}
	if ProtocolVersion() != 66051 {
		t.Error("ProtocolVersion() returned invalid version")
	}
}

func TestGenKey(t *testing.T) {
	a, err := GenKey(32)
	if err != nil {
		t.Error("GenKey failed:", err)
	}
	b, err := GenKey(32)
	if err != nil {
		t.Error("GenKey failed:", err)
	}
	if a.Cmp(b) == 0 {
		t.Error("GenKey returned the same key twice")
	}
}
