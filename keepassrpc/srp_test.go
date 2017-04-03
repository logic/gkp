package keepassrpc

import (
	"crypto/sha256"
	"fmt"
	"math/big"
	"os"
	"testing"
)

var mult = []byte("" +
	"\xb7\x86\x7f\x12\x99\xda\x8c\xc2\x4a\xb9" +
	"\x3e\x08\x98\x6e\xbc\x4d\x6a\x47\x8a\xd0" +
	"")
var k = new(big.Int).SetBytes(mult)

var username = "username"

var salt = "salt"
var password = "password"
var snp = "13601bda4ea78e55a07b98866d2be6be0744e3866f13c00c811cab608a28f322"
var pkey = new(big.Int)

var value = big.NewInt(1)
var A = new(big.Int).Exp(Generator, value, Prime)
var server = big.NewInt(2)
var pubs = "785f3ec7eb32f30b90cd0fcf3657d388b5ff4297f2f9716ff66e9b69c05ddd09"

var testClient *SRPContext

func TestMultiplier(t *testing.T) {
	p := testClient.Multiplier()
	if k.Cmp(p) != 0 {
		t.Error("Multipler() did not return standard result")
	}
}

func TestSetServer(t *testing.T) {
	if testClient.SetServer(server) != nil {
		t.Error("SetServer() failed")
	}
	if server.Cmp(testClient.Server) != 0 {
		t.Error("SetServer() set wrong value")
	}
}

func TestSetServerPrime(t *testing.T) {
	c := &SRPContext{
		Private: new(big.Int).Set(value),
		Public:  new(big.Int).Exp(Generator, value, Prime),
	}
	B := new(big.Int).Set(Prime)
	if c.SetServer(B) == nil {
		t.Error("SetServer() didn't catch prime match")
	}
}

func TestPrivateKey(t *testing.T) {
	if testClient.privateKey().Text(16) != snp {
		t.Error("privateKey() didn't match input")
	}
}

func TestScramblingParameter(t *testing.T) {
	if testClient.scramblingParameter().Text(16) != pubs {
		t.Error("scramblingParameter() didn't match input")
	}
}

func TestPremasterSecret(t *testing.T) {
	x, _ := new(big.Int).SetString(snp, 16)
	u, _ := new(big.Int).SetString(pubs, 16)
	gxN := new(big.Int).Exp(Generator, x, Prime)
	kgxN := new(big.Int).Exp(Generator, x, Prime).Mul(k, gxN)
	ux := new(big.Int).Mul(u, x)
	aux := new(big.Int).Add(value, ux)
	tmp := new(big.Int).Sub(server, kgxN)
	S := new(big.Int).Exp(tmp, aux, Prime)

	if S.Cmp(testClient.premasterSecret()) != 0 {
		t.Error("premasterSecret() didn't match input")
	}
}

func TestVerifier(t *testing.T) {
	x := new(big.Int).Exp(Generator, pkey, Prime)
	x1, err := testClient.Verifier()
	if err != nil {
		t.Error("Verifier() returned an error")
	}
	if x.Cmp(x1) != 0 {
		t.Error("Verifier() didn't match input")
	}
}

func TestVerifierFailure(t *testing.T) {
	c := &SRPContext{
		Private: new(big.Int).Set(value),
		Public:  new(big.Int).Exp(Generator, value, Prime),
	}
	if _, err := c.Verifier(); err == nil {
		t.Error("Verifier() failed to catch premature generation")
	}
}

func TestEvidence(t *testing.T) {
	S := testClient.premasterSecret()

	h := sha256.New()
	fmt.Fprintf(h, "%X%X%X", A, server, S)
	M := new(big.Int).SetBytes(h.Sum(nil))

	Mx, err := testClient.Evidence()
	if err != nil {
		t.Error("Evidence() returned an error")
	}
	if M.Cmp(Mx) != 0 {
		t.Error("Evidence() didn't match input")
	}
}

func TestEvidenceFailure(t *testing.T) {
	c := &SRPContext{
		Private: new(big.Int).Set(value),
		Public:  new(big.Int).Exp(Generator, value, Prime),
	}
	if _, err := c.Evidence(); err == nil {
		t.Error("Evidence() failed to catch premature generation")
	}
}

func TestServerEvidence(t *testing.T) {
	S := testClient.premasterSecret()
	M, err := testClient.Evidence()
	if err != nil {
		t.Error("Evidence() returned an error")
	}

	h := sha256.New()
	fmt.Fprintf(h, "%X%x%X", A, M, S)
	M1 := new(big.Int).SetBytes(h.Sum(nil))

	M1x, err := testClient.ServerEvidence()
	if err != nil {
		t.Error("ServerEvidence() returned an error")
	}
	if M1.Cmp(M1x) != 0 {
		t.Error("ServerEvidence() didn't match input")
	}
}

func TestServerEvidenceFailure(t *testing.T) {
	c := &SRPContext{
		Private: new(big.Int).Set(value),
		Public:  new(big.Int).Exp(Generator, value, Prime),
	}
	if _, err := c.ServerEvidence(); err == nil {
		t.Error("ServerEvidence() failed to catch premature generation")
	}
}

func TestMain(m *testing.M) {
	testClient = &SRPContext{
		Private:  new(big.Int).Set(value),
		Public:   new(big.Int).Exp(Generator, value, Prime),
		Salt:     salt,
		Password: password,
	}
	testClient.SetServer(server)

	var ok bool
	if pkey, ok = new(big.Int).SetString(snp, 16); !ok {
		panic("Conversion of snp failed")
	}

	os.Exit(m.Run())
}
