package network

import (
	"bytes"
	"encoding/binary"
	"errors"
)

type ConnHeaderCodec struct{}

func (codec *ConnHeaderCodec) Encode(header interface{}) ([]byte, error) {
	connHeader, ok := header.(*ConnHeader)
	if !ok {
		return nil, errors.New("invalid header type")
	}

	// The length of a UUID string is 36 characters (32 alphanumerics and 4 '-').
	// Combined with the 8-byte timestamp, the total buffer length is 36 + 8.
	buf := new(bytes.Buffer)

	// Write timestamp (8 bytes)
	binary.Write(buf, binary.BigEndian, connHeader.Timestamp)

	return buf.Bytes(), nil
}

func (codec *ConnHeaderCodec) Decode(data []byte) (interface{}, error) {
	// Read timestamp
	timestamp := int64(binary.BigEndian.Uint64(data[:8]))

	return &ConnHeader{
		Timestamp: timestamp,
	}, nil
}

type ConnAckHeaderCodec struct{}

func (codec *ConnAckHeaderCodec) Encode(header interface{}) ([]byte, error) {
	connAckHeader, ok := header.(*ConnAckHeader)
	if !ok {
		return nil, errors.New("invalid header type")
	}

	// The length of a UUID string is 36 characters (32 alphanumerics and 4 '-').
	// Combined with the 8-byte timestamp, the total buffer length is 36 + 8.
	buf := new(bytes.Buffer)

	// Write timestamp (8 bytes)
	binary.Write(buf, binary.BigEndian, connAckHeader.Timestamp)
	// Write UUID string (no fixed byte length)
	buf.WriteString(connAckHeader.Id)

	return buf.Bytes(), nil
}

func (codec *ConnAckHeaderCodec) Decode(data []byte) (interface{}, error) {
	// Read timestamp
	timestamp := int64(binary.BigEndian.Uint64(data[:8]))
	// Read UUID
	// The UUID is the rest of the buffer after the timestamp.
	id := string(data[8:])

	return &ConnAckHeader{
		Id:        id,
		Timestamp: timestamp,
	}, nil
}

type ConnHeader struct {
	// Timestamp
	Timestamp int64
}

type ConnAckHeader struct {
	// Client's Connection ID，定义更新使用UUID
	Id        string
	Timestamp int64
}
