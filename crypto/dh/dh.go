package dh

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"math/big"
)

// 定义基本参数
var (
	prime, _  = new(big.Int).SetString("FFFFFFFFFFFFFFFFC90FDAA22168C234C4C6628B80DC1CD129024E088A67CC74020BBEA63B139B22514A08798E3404DDEF9519B3CD3A431B302B0A6D7979FB1", 16)
	generator = big.NewInt(2)
)

// 生成随机私钥
func GenPrivateKey() *big.Int {
	privateKey, _ := rand.Int(rand.Reader, prime)
	return privateKey
}

// 通过私钥生成共享密钥
func GenSharedSecret(privateKey *big.Int) (*big.Int, error) {
	if privateKey == nil {
		return nil, errors.New("private key cannot be nil")
	}
	return new(big.Int).Exp(generator, privateKey, prime), nil
}

// 通过私钥和对方公钥计算会话密钥
func ComputeSecret(theirPublic, myPrivate *big.Int) ([]byte, error) {
	if theirPublic == nil || myPrivate == nil {
		return nil, errors.New("public key or private key is nil or both are nil")
	}
	secretInt := new(big.Int).Exp(theirPublic, myPrivate, prime)
	secretBytes := secretInt.Bytes()

	mac := hmac.New(sha256.New, secretBytes)
	mac.Write(secretBytes)

	return mac.Sum(nil), nil
}
