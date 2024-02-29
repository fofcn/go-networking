package aes

import (
	"crypto/aes"
	"crypto/cipher"
)

// var aesCipherTable map[string]interface{} = make(map[string]interface{})

func AesEncrypt(plainText []byte, key []byte, nonce []byte) ([]byte, error) {
	cipherBlock, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	cipherGCM, err := cipher.NewGCM(cipherBlock)
	if err != nil {
		return nil, err
	}

	return cipherGCM.Seal(nil, nonce, plainText, nil), nil
}

func AesDecrypt(cipherText []byte, nonce []byte, key []byte) ([]byte, error) {
	cipherBlock, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	cipherGCM, err := cipher.NewGCM(cipherBlock)
	if err != nil {
		return nil, err
	}

	decrypted, err := cipherGCM.Open(nil, nonce, cipherText, nil)
	if err != nil {
		return nil, err
	}

	return decrypted, nil
}
