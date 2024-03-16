package network_test

import (
	"go-networking/network"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCloseHeaderCodec_Encode_ShouldReturnEncodedBytes_WhenGivenValidHeader(t *testing.T) {
	header := network.CloseHeader{
		Id:     "12345",
		Reason: "Test reason",
	}
	codec := network.CloseHeaderCodec{}
	encoded, err := codec.Encode(&header)
	assert.NoError(t, err)
	assert.NotNil(t, encoded)
}

func TestCloseHeaderCodec_Encode_ShouldReturnError_WhenGivenInvalidType(t *testing.T) {
	codec := network.CloseHeaderCodec{}
	_, err := codec.Encode("invalid type")
	assert.Error(t, err)
}

func TestCloseHeaderCodec_Decode_ShouldReturnCloseHeader_WhenGivenEncodedBytes(t *testing.T) {
	// Use the output from the previous test or manually construct an encoded example
	encoded := []byte{0, 5, '1', '2', '3', '4', '5', 0, 11, 'T', 'e', 's', 't', ' ', 'r', 'e', 'a', 's', 'o', 'n'}
	codec := network.CloseHeaderCodec{}
	decoded, err := codec.Decode(encoded)
	assert.NoError(t, err)

	header, ok := decoded.(*network.CloseHeader)
	assert.True(t, ok)
	assert.Equal(t, "12345", header.Id)
	assert.Equal(t, "Test reason", header.Reason)
}

func TestCloseHeaderCodec_Decode_ShouldReturnError_WhenGivenInvalidData(t *testing.T) {
	codec := network.CloseHeaderCodec{}
	_, err := codec.Decode([]byte{0})
	assert.Error(t, err)
}

func TestCloseAckHeaderCodec_Encode_ShouldReturnEncodedBytes_WhenGivenValidHeader(t *testing.T) {
	header := network.CloseAckHeader{
		StatusCode: 200,
		Details:    "OK",
	}
	codec := network.CloseAckHeaderCodec{}
	encoded, err := codec.Encode(&header)
	assert.NoError(t, err)
	assert.NotNil(t, encoded)
}

func TestCloseAckHeaderCodec_Encode_ShouldReturnError_WhenGivenInvalidType(t *testing.T) {
	codec := network.CloseAckHeaderCodec{}
	_, err := codec.Encode("invalid type")
	assert.Error(t, err)
}

func TestCloseAckHeaderCodec_Decode_ShouldReturnCloseAckHeader_WhenGivenEncodedBytes(t *testing.T) {
	// Use the output from the previous test or manually construct an encoded example
	encoded := []byte{0, 200, 0, 2, 'O', 'K'}
	codec := network.CloseAckHeaderCodec{}
	decoded, err := codec.Decode(encoded)
	assert.NoError(t, err)

	header, ok := decoded.(*network.CloseAckHeader)
	assert.True(t, ok)
	assert.Equal(t, uint16(200), header.StatusCode)
	assert.Equal(t, "OK", header.Details)
}

func TestCloseAckHeaderCodec_Decode_ShouldReturnError_WhenGivenInvalidData(t *testing.T) {
	codec := network.CloseAckHeaderCodec{}
	_, err := codec.Decode([]byte{0})
	assert.Error(t, err)
}
