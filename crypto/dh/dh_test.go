package dh_test

import (
	"fmt"
	"go-networking/crypto/dh"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenDHKeyShouldReturnSuccesWhenTwoKeysGeneratedSuccess(t *testing.T) {
	//每个通讯方生成各自的私钥
	alicePrivateKey := dh.GenPrivateKey()
	bobPrivateKey := dh.GenPrivateKey()

	//然后生成对应的公钥
	alicePublicKey, err := dh.GenSharedSecret(alicePrivateKey)
	assert.Nil(t, err)

	bobPublicKey, err := dh.GenSharedSecret(bobPrivateKey)
	assert.Nil(t, err)

	//通过私钥和对方公钥计算出会话密钥
	aliceSessionKey, err := dh.ComputeSecret(bobPublicKey, alicePrivateKey)
	assert.Nil(t, err)

	bobSessionKey, err := dh.ComputeSecret(alicePublicKey, bobPrivateKey)
	assert.Nil(t, err)

	fmt.Printf("Alice session key: %s\n", aliceSessionKey)
	fmt.Printf("Bob session key: %s\n", bobSessionKey)

	assert.Equal(t, aliceSessionKey, bobSessionKey)
}

func TestGenDHKeyShouldFailedWhenOneKeyIsNotGenerated(t *testing.T) {
	// given
	bobPrivateKey := dh.GenPrivateKey()

	bobPublicKey, err := dh.GenSharedSecret(bobPrivateKey)
	assert.Nil(t, err)

	_, err = dh.ComputeSecret(bobPublicKey, nil)
	assert.NotNil(t, err)
}
