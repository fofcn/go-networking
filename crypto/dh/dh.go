package dh

import (
	"crypto/rand"
	"crypto/sha256"
	"math/big"
)

// Generate DH Key Pair
func GenDHKP(p, g *big.Int) (*big.Int, *big.Int, error) {
	privateKey, err := rand.Int(rand.Reader, p)
	if err != nil {
		return nil, nil, err
	}
	publicKey := new(big.Int).Exp(g, privateKey, p)
	return privateKey, publicKey, nil
}

func FastGenDHKP() (*big.Int, *big.Int, error) {
	// mod P
	p := new(big.Int).SetInt64(23)
	// base G
	g := new(big.Int).SetInt64(5)
	return GenDHKP(p, g)
}

// Generate DH shared key
func GenDHSharedKey(otherPublicKey, myPrivateKey, p *big.Int) *big.Int {
	// sharedKey = (otherPublicKey^myPrivateKey) mod p
	return new(big.Int).Exp(otherPublicKey, myPrivateKey, p)
}

func FastGenDHSharedKey(otherPublicKey, myPrivateKey *big.Int) *big.Int {
	p := new(big.Int).SetInt64(23) // mod P
	return GenDHSharedKey(otherPublicKey, myPrivateKey, p)
}

func GenAESKeyFromDHKey(dhKey *big.Int) []byte {
	dhSharedKeyBytes := dhKey.Bytes()
	hashedKey := sha256.Sum256(dhSharedKeyBytes)
	return hashedKey[:]
}
