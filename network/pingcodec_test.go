package network_test

import (
	"go-networking/network"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const uuidRawLength = 32

func TestPingHeaderCodec_EncodeDecode_ShouldEncodeAndDecodeCorrectly_WhenGivenValidHeader(t *testing.T) {
	// Setup example PingHeader
	pingHeader := &network.PingHeader{
		Timestamp: time.Now().Unix(),
		Id:        "123e4567e89b12d3a456426614174000",
	}

	// Initialize codec
	codec := network.PingHeaderCodec{}

	// Encode
	encoded, err := codec.Encode(pingHeader)
	assert.NoError(t, err)
	assert.Len(t, encoded, 8+uuidRawLength)

	// Decode
	decoded, err := codec.Decode(encoded)
	assert.NoError(t, err)

	decodedPingHeader, ok := decoded.(*network.PingHeader)
	assert.True(t, ok)
	assert.Equal(t, pingHeader.Timestamp, decodedPingHeader.Timestamp)
	assert.Equal(t, pingHeader.Id, decodedPingHeader.Id)
}

func TestPongHeaderCodec_EncodeDecode_ShouldEncodeAndDecodeCorrectly_WhenGivenValidHeader(t *testing.T) {
	// Setup example PongHeader
	pongHeader := &network.PongHeader{
		Timestamp: time.Now().Unix(),
	}

	// Initialize codec
	codec := network.PongHeaderCodec{}

	// Encode
	encoded, err := codec.Encode(pongHeader)
	assert.NoError(t, err)
	assert.Len(t, encoded, 8)

	// Decode
	decoded, err := codec.Decode(encoded)
	assert.NoError(t, err)

	decodedPongHeader, ok := decoded.(*network.PongHeader)
	assert.True(t, ok)
	assert.Equal(t, pongHeader.Timestamp, decodedPongHeader.Timestamp)
}
