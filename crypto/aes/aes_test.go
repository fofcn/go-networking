package aes_test

import (
	"crypto/rand"
	"fmt"
	"go-networking/crypto/aes"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestAesEncryptShouldReturnErrorWhenKeyIsInvalid checks that AesEncrypt returns an error when the key is invalid.
func TestAesEncryptShouldReturnErrorWhenKeyIsInvalid(t *testing.T) {
	key := []byte{0} // Invalid key
	plaintext := []byte("hello, world")
	nonce := make([]byte, 12)

	_, err := aes.AesEncrypt(plaintext, key, nonce)

	assert.Error(t, err)
}

// TestAesEncryptAndAesDecryptShouldReturnSamePlaintext checks that AesEncrypt and AesDecrypt returns the same plaintext.
func TestAesEncryptAndAesDecryptShouldReturnSamePlaintext(t *testing.T) {
	key := make([]byte, 32)
	_, _ = rand.Read(key) // or replace with your actual key

	plaintext := []byte("hello, world")
	nonce := make([]byte, 12)
	ciphertext, err := aes.AesEncrypt(plaintext, key, nonce)
	assert.NoError(t, err, "AesEncrypt should not return an error")

	decrypted, err := aes.AesDecrypt(ciphertext, nonce, key)
	assert.NoError(t, err, "AesDecrypt should not return an error")

	assert.Equal(t, plaintext, decrypted, "The decrypted text should equal the plaintext")
}

// TestAesDecryptShouldReturnErrorWhenNonceIsIncorrect checks that AesDecrypt returns an error when nonce is incorrect.
func TestAesDecryptShouldReturnErrorWhenNonceIsIncorrect(t *testing.T) {
	key := make([]byte, 32)
	_, _ = rand.Read(key) // or replace with your actual key

	plaintext := []byte("hello, world")
	fmt.Printf("%x\n", plaintext)
	correctnonce := []byte("987654321012")

	ciphertext, _ := aes.AesEncrypt(plaintext, key, correctnonce)

	nonce := []byte("123456789012") // Incorrect nonce

	incorrectPlaintext, err := aes.AesDecrypt(ciphertext, nonce, key)
	fmt.Printf("%x %s\n", incorrectPlaintext, err)

	assert.NotEqual(t, plaintext, incorrectPlaintext)
}

// TestAesDecryptShouldReturnErrorWhenKeyIsInvalid 测试当密钥无效时 AesDecrypt 应返回错误
func TestAesDecryptShouldReturnErrorWhenKeyIsInvalid(t *testing.T) {
	key := []byte{0} // Invalid key
	ciphertext := []byte("Some encrypted text")
	nonce := make([]byte, 12)

	_, err := aes.AesDecrypt(ciphertext, nonce, key)
	assert.Error(t, err)
}

// TestAesDecryptShouldReturnErrorWhenCiphertextIsInvalid 测试当密文无效时 AesDecrypt 应返回错误
func TestAesDecryptShouldReturnErrorWhenCiphertextIsInvalid(t *testing.T) {
	key := make([]byte, 32)
	_, _ = rand.Read(key) // or replace with your actual key

	ciphertext := []byte("Invalid ciphertext")
	nonce := make([]byte, 12)

	_, err := aes.AesDecrypt(ciphertext, nonce, key)
	assert.Error(t, err)
}

// TestAesDecryptShouldReturnPlaintextWhenInputsAreValid 测试当所有输入都有效时 AesDecrypt 应返回原始明文
func TestAesDecryptShouldReturnPlaintextWhenInputsAreValid(t *testing.T) {
	key := make([]byte, 32)
	_, _ = rand.Read(key) // replace with your actual key
	nonce := make([]byte, 12)

	plaintext := []byte("hello, world")
	ciphertext, _ := aes.AesEncrypt(plaintext, key, nonce)

	decrypted, err := aes.AesDecrypt(ciphertext, nonce, key)

	assert.NoError(t, err)
	assert.Equal(t, plaintext, decrypted)
}
