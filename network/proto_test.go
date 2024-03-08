package network_test

import (
	"bytes"
	"encoding/binary"
	"errors"
	"go-networking/network"
	"testing"
)

const (
	CommandA network.CommandType = iota
)

type Conn struct {
	KeyLen uint32
	Key    string
}

type ConnCodec struct {
}

func (codec ConnCodec) Encode(header interface{}) ([]byte, error) {
	if conn, ok := header.(Conn); ok {
		buf := new(bytes.Buffer)
		buf.Write(network.EncodeInteger(uint64(conn.KeyLen)))
		buf.Write([]byte(conn.Key))
		return buf.Bytes(), nil
	} else {
		// todo
		return nil, errors.New("error occured when try to encode, invalid type of Conn")
	}
}
func (codec ConnCodec) Decode(data []byte) (interface{}, error) {
	buf := bytes.NewReader(data)
	keyLen, err := binary.ReadUvarint(buf)
	if err != nil {
		return nil, err
	}

	conn := &Conn{}
	conn.KeyLen = uint32(keyLen)

	keyData := make([]byte, keyLen)

	if _, err := buf.Read(keyData); err != nil {
		return nil, errors.New("failed to read payload")
	}

	conn.Key = string(keyData)
	return conn, nil
}

func TestEncodeShouldReturnBytesWhenEncodeSuccess(t *testing.T) {
	network.AddCodec(CommandA, &ConnCodec{})

	givenConn := Conn{
		KeyLen: uint32(len("ABC")),
		Key:    "ABC",
	}
	proto := &network.Frame{
		Version:  1,
		CmdType:  CommandA,
		Sequence: 1,
		Header:   givenConn,
		Payload:  []byte{0x04, 0x05, 0x06},
	}

	data, err := network.Encode(proto)
	if err != nil {
		t.Fatalf("Expected nil error, got %v", err)
	}

	expectedData := []byte{0x01, 0x00, 0x01, 0x04, 0x03, 0x41, 0x42, 0x43, 0x04, 0x05, 0x06}
	if !bytes.Equal(data, expectedData) {
		t.Errorf("Expected %v, got %v", expectedData, data)
	}
}

func TestEncodeShouldReturnErrorWhenHeaderEncodeFails(t *testing.T) {
	network.AddCodec(CommandA, &ConnCodec{})
	proto := &network.Frame{}
	_, err := network.Encode(proto)
	if err == nil {
		t.Error("Expected error, got nil")
	}
}

func TestDecodeShouldReturnFrameWhenDecodeSuccess(t *testing.T) {
	network.AddCodec(CommandA, &ConnCodec{})

	frame := []byte{0x01, 0x00, 0x01, 0x03, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06}
	proto, err := network.Decode(frame)
	if err != nil {
		t.Fatalf("Expected nil error, got %v", err)
	}

	if proto.Version != 0x0001 {
		t.Errorf("Expected %v, got %v", 0x0001, proto.Version)
	}

	if proto.CmdType != CommandA {
		t.Errorf("Expected %v, got %v", CommandA, proto.CmdType)
	}

	expectedPayload := []byte{0x04, 0x05, 0x06}
	if !bytes.Equal(proto.Payload, expectedPayload) {
		t.Errorf("Expected %v, got %v", expectedPayload, proto.Payload)
	}
}

func TestDecodeShouldReturnErrorWhenFrameTooShort(t *testing.T) {
	shortFrame := []byte{0x01, 0x02}

	_, err := network.Decode(shortFrame)
	if err == nil {
		t.Error("Expected error, got nil")
	}
}

func TestDecodeShouldReturnErrorWhenSubheaderLengthReadFails(t *testing.T) {
	// 创建一个只包含版本、命令和长度字段的帧，但长度字段的长度超出了剩余的帧大小
	frame := []byte{0x00, 0x01, 0x00, 0x01, 0x00, 0x04}
	_, err := network.Decode(frame)
	if err == nil {
		t.Error("Expected error, got nil")
	}
}
