package network

import (
	"encoding/binary"
	"errors"
)

type PingHeader struct {
	Timestamp int64 // Timestamp
	// Client's Connection ID，定义更新使用UUID
	Id string
}

type PongHeader struct {
	Timestamp int64 // Timestamp
}

const uuidRawLength = 32 // Length of UUID without hyphens

type PingHeaderCodec struct{}

func (codec *PingHeaderCodec) Encode(header interface{}) ([]byte, error) {
	pingHeader, ok := header.(*PingHeader)
	if !ok {
		return nil, errors.New("invalid header type for PING")
	}

	// Allocate a buffer for the timestamp (8 bytes) + UUID (32 bytes)
	buf := make([]byte, 8+uuidRawLength)

	// Write timestamp (8 bytes)
	binary.BigEndian.PutUint64(buf[0:8], uint64(pingHeader.Timestamp))

	// Write UUID (32 bytes)
	copy(buf[8:], pingHeader.Id)

	return buf, nil
}

func (codec *PingHeaderCodec) Decode(data []byte) (interface{}, error) {
	if len(data) < (8 + uuidRawLength) {
		return nil, errors.New("data too short for decoding PING header")
	}

	// Read timestamp (8 bytes)
	timestamp := int64(binary.BigEndian.Uint64(data[:8]))

	// Read UUID (32 bytes)
	id := string(data[8 : 8+uuidRawLength])

	return &PingHeader{
		Timestamp: timestamp,
		Id:        id,
	}, nil
}

type PongHeaderCodec struct{}

func (codec *PongHeaderCodec) Encode(header interface{}) ([]byte, error) {
	pongHeader, ok := header.(*PongHeader)
	if !ok {
		return nil, errors.New("invalid header type for PONG")
	}

	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(pongHeader.Timestamp))
	return buf, nil
}

func (codec *PongHeaderCodec) Decode(data []byte) (interface{}, error) {
	if len(data) < 8 {
		return nil, errors.New("data too short for decoding PONG header")
	}

	timestamp := int64(binary.BigEndian.Uint64(data[:8]))
	return &PongHeader{Timestamp: timestamp}, nil
}
