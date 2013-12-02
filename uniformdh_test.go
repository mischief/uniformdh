package uniformdh

import (
	"testing"
  "bytes"
)

func TestUniformDH(t *testing.T) {
	alice := New()
	bob := New()
	alicePub, bobPub := alice.Public(), bob.Public()
	aliceSecret, bobSecret := alice.Secret(bobPub), bob.Secret(alicePub)

	if !bytes.Equal(aliceSecret, bobSecret) {
		t.Fatalf("secret key differs")
		t.Logf("alice: %d bytes", len(aliceSecret))
		t.Logf("%X", aliceSecret)
		t.Logf("bob: %d bytes", len(bobSecret))
		t.Logf("%X", bobSecret)
	}
}
