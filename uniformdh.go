// Package uniformdh implements the UniformDH key exchange
// algorithm as used in the obfs3 tor pluggable transport.
//
// See https://gitweb.torproject.org/pluggable-transports/obfsproxy.git/blob_plain/HEAD:/doc/obfs3/obfs3-protocol-spec.txt
// for details.
package uniformdh

import (
	"crypto/rand"
	"math/big"
)

var (
	// generator of group 5
	g int64 = 2

	// byte size of group 5
	groupLen = 192

	// bitsize of group 5
	intSize = 1536

	// 1536 bit MODP group 5 from RFC 3526
	// equivalent to 2^1536 - 2^1472 - 1 + 2^64 * { [2^1406 pi] + 741804 }
	mod      big.Int
	modBytes = []byte{
		0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xC9, 0x0F, 0xDA, 0xA2,
		0x21, 0x68, 0xC2, 0x34, 0xC4, 0xC6, 0x62, 0x8B, 0x80, 0xDC, 0x1C, 0xD1,
		0x29, 0x02, 0x4E, 0x08, 0x8A, 0x67, 0xCC, 0x74, 0x02, 0x0B, 0xBE, 0xA6,
		0x3B, 0x13, 0x9B, 0x22, 0x51, 0x4A, 0x08, 0x79, 0x8E, 0x34, 0x04, 0xDD,
		0xEF, 0x95, 0x19, 0xB3, 0xCD, 0x3A, 0x43, 0x1B, 0x30, 0x2B, 0x0A, 0x6D,
		0xF2, 0x5F, 0x14, 0x37, 0x4F, 0xE1, 0x35, 0x6D, 0x6D, 0x51, 0xC2, 0x45,
		0xE4, 0x85, 0xB5, 0x76, 0x62, 0x5E, 0x7E, 0xC6, 0xF4, 0x4C, 0x42, 0xE9,
		0xA6, 0x37, 0xED, 0x6B, 0x0B, 0xFF, 0x5C, 0xB6, 0xF4, 0x06, 0xB7, 0xED,
		0xEE, 0x38, 0x6B, 0xFB, 0x5A, 0x89, 0x9F, 0xA5, 0xAE, 0x9F, 0x24, 0x11,
		0x7C, 0x4B, 0x1F, 0xE6, 0x49, 0x28, 0x66, 0x51, 0xEC, 0xE4, 0x5B, 0x3D,
		0xC2, 0x00, 0x7C, 0xB8, 0xA1, 0x63, 0xBF, 0x05, 0x98, 0xDA, 0x48, 0x36,
		0x1C, 0x55, 0xD3, 0x9A, 0x69, 0x16, 0x3F, 0xA8, 0xFD, 0x24, 0xCF, 0x5F,
		0x83, 0x65, 0x5D, 0x23, 0xDC, 0xA3, 0xAD, 0x96, 0x1C, 0x62, 0xF3, 0x56,
		0x20, 0x85, 0x52, 0xBB, 0x9E, 0xD5, 0x29, 0x07, 0x70, 0x96, 0x96, 0x6D,
		0x67, 0x0C, 0x35, 0x4E, 0x4A, 0xBC, 0x98, 0x04, 0xF1, 0x74, 0x6C, 0x08,
		0xCA, 0x23, 0x73, 0x27, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
	}
)

func init() {
	// setup 1536 bit MODP group
	mod.SetBytes(modBytes)
}

type UniformDH struct {
	priv         big.Int
	pubStr       []byte
	pub          big.Int
	sharedSecret big.Int
}

// Create a new UniformDH instance
func New() *UniformDH {
	udh := &UniformDH{}

	privStr := make([]byte, groupLen)
	// To pick a private UniformDH key, we pick a random 1536-bit number,
	// and make it even by setting its low bit to 0. Let x be that private
	// key, and X = g^x (mod p).
	rand.Read(privStr)
	udh.priv.SetBytes(privStr)

	// When someone sends her public key to the other party, she randomly
	// decides whether to send X or p-X. This makes the public key
	// negligibly different from a uniform 1536-bit string
	flip := new(big.Int).Mod(&udh.priv, big.NewInt(2))
	udh.priv.Sub(&udh.priv, flip)

	udh.pub.Exp(big.NewInt(g), &udh.priv, &mod)

	if flip.Uint64() == 1 {
		udh.pub.Sub(&mod, &udh.pub)
	}

	/// XXX: handle erroneous situations better
	if udh.priv.BitLen() > intSize {
		panic("int too large")
	}

	return udh
}

// Returns big-endian public key
func (udh *UniformDH) Public() *[192]byte {
	var buf [192]byte
	pubBytes := udh.pub.Bytes()

	copy(buf[groupLen-len(pubBytes):], pubBytes)

	return &buf
}

// Returns the shared secret, given the other party's public key
func (udh *UniformDH) Secret(theirPubBytes *[192]byte) (secret *[192]byte) {
	// When a party wants to calculate the shared secret, she
	// raises the foreign public key to her private key. Note that both
	// (p-Y)^x = Y^x (mod p) and (p-X)^y = X^y (mod p), since x and y are
	// even.
	theirPub := new(big.Int).SetBytes(theirPubBytes[:])
	udh.sharedSecret.Exp(theirPub, &udh.priv, &mod)

	sharedBytes := udh.sharedSecret.Bytes()

	var buf [192]byte

	copy(buf[groupLen-len(sharedBytes):], sharedBytes)

	return &buf
}
