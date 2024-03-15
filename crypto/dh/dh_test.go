package dh_test

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"go-networking/crypto/dh"
	"io"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenDHKPShouldReturnKeysWhenGivenPG(t *testing.T) {
	p := big.NewInt(23)
	g := big.NewInt(5)
	privateKey, publicKey, err := dh.GenDHKP(p, g)
	require.NoError(t, err)
	assert.NotNil(t, privateKey)
	assert.NotNil(t, publicKey)
}

func TestFastGenDHKPShouldReturnKeys(t *testing.T) {
	privateKey, publicKey, err := dh.FastGenDHKP()
	require.NoError(t, err)
	assert.NotNil(t, privateKey)
	assert.NotNil(t, publicKey)
}

func TestGenDHSharedKeyShouldReturnSharedKeyWhenGivenOtherPublicKeyMyPrivateKeyAndP(t *testing.T) {
	// using random key as other's public key and my private key
	otherPublicKey, _ := rand.Int(rand.Reader, big.NewInt(23))
	myPrivateKey, _ := rand.Int(rand.Reader, big.NewInt(23))
	p := big.NewInt(23)

	sharedKey := dh.GenDHSharedKey(otherPublicKey, myPrivateKey, p)
	assert.NotNil(t, sharedKey)
}

func TestFastGenDHSharedKeyShouldReturnSharedKeyWhenGivenOtherPublicKeyAndMyPrivateKey(t *testing.T) {
	// using random key as other's public key and my private key
	otherPublicKey, _ := rand.Int(rand.Reader, big.NewInt(23))
	myPrivateKey, _ := rand.Int(rand.Reader, big.NewInt(23))

	sharedKey := dh.FastGenDHSharedKey(otherPublicKey, myPrivateKey)
	assert.NotNil(t, sharedKey)
}

func TestGenAESKeyFromDHKeyShouldReturnAESKeyWhenGivenDHKey(t *testing.T) {
	dhKey := big.NewInt(23)
	aesKey := dh.GenAESKeyFromDHKey(dhKey)
	assert.NotNil(t, aesKey)
	assert.Equal(t, 32, len(aesKey)) // SHA-256 的输出长度是32字节
}

func TestAliceBobShouldGetSameSharedKey(t *testing.T) {
	// Alice and Bob generate their own private and public keys
	privateKeyAlice, publicKeyAlice, err := dh.FastGenDHKP()
	require.NoError(t, err)
	privateKeyBob, publicKeyBob, err := dh.FastGenDHKP()
	require.NoError(t, err)

	// Alice and Bob generate shared key from each other's public key and their own private key
	sharedKeyAlice := dh.FastGenDHSharedKey(publicKeyBob, privateKeyAlice)
	sharedKeyBob := dh.FastGenDHSharedKey(publicKeyAlice, privateKeyBob)

	assert.Equal(t, sharedKeyAlice, sharedKeyBob)
}

func TestAliceBobShouldGetSameCipherTextWhenGivenPlainText(t *testing.T) {
	// Alice and Bob generate their own private and public keys
	privateKeyAlice, publicKeyAlice, err := dh.FastGenDHKP()
	require.NoError(t, err)
	privateKeyBob, publicKeyBob, err := dh.FastGenDHKP()
	require.NoError(t, err)

	// Alice and Bob generate shared key from each other's public key and their own private key
	sharedKeyAlice := dh.FastGenDHSharedKey(publicKeyBob, privateKeyAlice)
	sharedKeyBob := dh.FastGenDHSharedKey(publicKeyAlice, privateKeyBob)

	// Alice and Bob generate AES key from shared key
	aesKeyAlice := dh.GenAESKeyFromDHKey(sharedKeyAlice)
	aesKeyBob := dh.GenAESKeyFromDHKey(sharedKeyBob)

	// Plaintext
	plainText := []byte("Hello, World!")

	blockAlice, err := aes.NewCipher(aesKeyAlice)
	require.NoError(t, err)
	blockBob, err := aes.NewCipher(aesKeyBob)
	require.NoError(t, err)

	cipherTextAlice := make([]byte, aes.BlockSize+len(plainText))
	cipherTextBob := make([]byte, aes.BlockSize+len(plainText))

	ivAlice := cipherTextAlice[:aes.BlockSize]
	ivBob := cipherTextBob[:aes.BlockSize]
	_, err = io.ReadFull(rand.Reader, ivAlice)
	require.NoError(t, err)
	copy(ivBob, ivAlice)
	require.NoError(t, err)

	streamAlice := cipher.NewCFBEncrypter(blockAlice, ivAlice)
	streamAlice.XORKeyStream(cipherTextAlice[aes.BlockSize:], plainText)

	streamBob := cipher.NewCFBEncrypter(blockBob, ivBob)
	streamBob.XORKeyStream(cipherTextBob[aes.BlockSize:], plainText)

	assert.Equal(t, cipherTextAlice, cipherTextBob)
}

func TestAliceBobShouldEncryptAndDecryptSuccessfullyGivenPlainText(t *testing.T) {
	// Alice and Bob generate their own private and public keys
	privateKeyAlice, publicKeyAlice, err := dh.FastGenDHKP()
	require.NoError(t, err)
	privateKeyBob, publicKeyBob, err := dh.FastGenDHKP()
	require.NoError(t, err)

	// Alice and Bob generate shared key from each other's public key and their own private key
	sharedKeyAlice := dh.FastGenDHSharedKey(publicKeyBob, privateKeyAlice)
	sharedKeyBob := dh.FastGenDHSharedKey(publicKeyAlice, privateKeyBob)

	// Alice and Bob generate AES key from shared key
	aesKeyAlice := dh.GenAESKeyFromDHKey(sharedKeyAlice)
	aesKeyBob := dh.GenAESKeyFromDHKey(sharedKeyBob)

	// Plaintext
	plainText := []byte("Hello, World!")

	blockAlice, err := aes.NewCipher(aesKeyAlice)
	require.NoError(t, err)
	blockBob, err := aes.NewCipher(aesKeyBob)
	require.NoError(t, err)

	cipherTextAlice := make([]byte, aes.BlockSize+len(plainText))
	cipherTextBob := make([]byte, aes.BlockSize+len(plainText))

	ivAlice := cipherTextAlice[:aes.BlockSize]
	ivBob := cipherTextBob[:aes.BlockSize]
	_, err = io.ReadFull(rand.Reader, ivAlice)
	require.NoError(t, err)
	copy(ivBob, ivAlice)

	streamAlice := cipher.NewCFBEncrypter(blockAlice, ivAlice)
	streamAlice.XORKeyStream(cipherTextAlice[aes.BlockSize:], plainText)

	streamBob := cipher.NewCFBEncrypter(blockBob, ivBob)
	streamBob.XORKeyStream(cipherTextBob[aes.BlockSize:], plainText)

	// Test encryption
	assert.Equal(t, cipherTextAlice, cipherTextBob)

	// Test decryption
	decryptedTextAlice := make([]byte, len(plainText))
	decryptedTextBob := make([]byte, len(plainText))

	streamAlice = cipher.NewCFBDecrypter(blockAlice, ivAlice)
	streamAlice.XORKeyStream(decryptedTextAlice, cipherTextAlice[aes.BlockSize:])

	streamBob = cipher.NewCFBDecrypter(blockBob, ivBob)
	streamBob.XORKeyStream(decryptedTextBob, cipherTextBob[aes.BlockSize:])

	assert.Equal(t, plainText, decryptedTextAlice)
	assert.Equal(t, plainText, decryptedTextBob)
}
