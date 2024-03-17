package codec

import (
	"bytes"
	"encoding/binary"
	"errors"
)

type CloseHeader struct {
	Id     string
	Reason string
}

type CloseAckHeader struct {
	StatusCode uint16
	Details    string
}

type CloseHeaderCodec struct{}

func (codec *CloseHeaderCodec) Encode(header interface{}) ([]byte, error) {
	closeHeader, ok := header.(*CloseHeader)
	if !ok {
		return nil, errors.New("invalid header type for CLOSE")
	}

	buf := new(bytes.Buffer)

	// Write ID and Reason strings with their lengths as prefixes
	WriteLvString(buf, closeHeader.Id)
	WriteLvString(buf, closeHeader.Reason)

	return buf.Bytes(), nil
}

func (codec *CloseHeaderCodec) Decode(data []byte) (interface{}, error) {
	reader := bytes.NewReader(data)

	// Read ID and Reason strings
	id, err := ReadLVString(reader)
	if err != nil {
		return nil, err
	}
	reason, err := ReadLVString(reader)
	if err != nil {
		return nil, err
	}

	return &CloseHeader{
		Id:     id,
		Reason: reason,
	}, nil
}

type CloseAckHeaderCodec struct{}

func (codec *CloseAckHeaderCodec) Encode(header interface{}) ([]byte, error) {
	closeAckHeader, ok := header.(*CloseAckHeader)
	if !ok {
		return nil, errors.New("invalid header type for CLOSEACK")
	}

	buf := new(bytes.Buffer)

	// Write StatusCode (2 bytes)
	if err := binary.Write(buf, binary.BigEndian, closeAckHeader.StatusCode); err != nil {
		return nil, err
	}

	// Write Details string with its length as prefix
	WriteLvString(buf, closeAckHeader.Details)

	return buf.Bytes(), nil
}

func (codec *CloseAckHeaderCodec) Decode(data []byte) (interface{}, error) {
	reader := bytes.NewReader(data)

	// Read StatusCode (2 bytes)
	var statusCode uint16
	if err := binary.Read(reader, binary.BigEndian, &statusCode); err != nil {
		return nil, err
	}

	// Read Details string
	details, err := ReadLVString(reader)
	if err != nil {
		return nil, err
	}

	return &CloseAckHeader{
		StatusCode: statusCode,
		Details:    details,
	}, nil
}

// Helper functions to write and read prefixed strings to buffered data
func WriteLvString(buf *bytes.Buffer, s string) {
	length := uint16(len(s))
	binary.Write(buf, binary.BigEndian, length)
	buf.WriteString(s)
}

func ReadLVString(reader *bytes.Reader) (string, error) {
	var length uint16
	if err := binary.Read(reader, binary.BigEndian, &length); err != nil {
		return "", err
	}
	strBytes := make([]byte, length)
	if _, err := reader.Read(strBytes); err != nil {
		return "", err
	}
	return string(strBytes), nil
}
