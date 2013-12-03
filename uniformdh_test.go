package uniformdh

import (
	"bytes"
	"testing"
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

func BenchmarkUniformDH(b *testing.B) {
	for i := 0; i < b.N; i++ {
		alice := New()
		bob := New()
		alicePub, bobPub := alice.Public(), bob.Public()
		aliceSecret, bobSecret := alice.Secret(bobPub), bob.Secret(alicePub)

		if !bytes.Equal(aliceSecret, bobSecret) {
			b.Fatalf("secret key differs")
			b.Logf("alice: %d bytes", len(aliceSecret))
			b.Logf("%X", aliceSecret)
			b.Logf("bob: %d bytes", len(bobSecret))
			b.Logf("%X", bobSecret)
		}
	}
}
