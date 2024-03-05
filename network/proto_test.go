package network_test

import (
	"bytes"
	"errors"
	"go-networking/network"
	"testing"
)

const (
	CommandA network.CommandType = iota
)

type MockVarintHeader struct {
	EncodeFunc func() ([]byte, error)
	DecodeFunc func([]byte) error
}

func (m *MockVarintHeader) Encode() ([]byte, error) {
	return m.EncodeFunc()
}

func (m *MockVarintHeader) Decode(b []byte) error {
	return m.DecodeFunc(b)
}

func TestEncodeShouldReturnBytesWhenEncodeSuccess(t *testing.T) {
	mockHeader := &MockVarintHeader{
		EncodeFunc: func() ([]byte, error) {
			// 模拟成功编码，返回预设的字节数组
			return []byte{0x01, 0x02, 0x03}, nil
		},
	}

	proto := &network.Proto{
		Version:      0x0001,
		CmdType:      CommandA,
		VarintHeader: mockHeader,
		Payload:      []byte{0x04, 0x05, 0x06},
	}

	data, err := network.Encode(proto)
	if err != nil {
		t.Fatalf("Expected nil error, got %v", err)
	}

	expectedData := []byte{0x00, 0x01, 0x00, 0x00, 0x00, 0x03, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06}
	if !bytes.Equal(data, expectedData) {
		t.Errorf("Expected %v, got %v", expectedData, data)
	}
}

func TestEncodeShouldReturnErrorWhenHeaderEncodeFails(t *testing.T) {
	mockHeader := &MockVarintHeader{
		EncodeFunc: func() ([]byte, error) {
			// 模拟编码失败
			return nil, errors.New("Failed to encode")
		},
		DecodeFunc: func([]byte) error { return nil },
	}
	proto := &network.Proto{VarintHeader: mockHeader}

	_, err := network.Encode(proto)
	if err == nil {
		t.Error("Expected error, got nil")
	}
}

func TestDecodeShouldReturnProtoWhenDecodeSuccess(t *testing.T) {
	mockHeader := &MockVarintHeader{
		DecodeFunc: func([]byte) error {
			// 模拟成功解码
			return nil
		},
	}

	cmdFactory.commandToCodec[CommandA] = mockHeader

	frame := []byte{0x00, 0x01, 0x00, 0x00, 0x00, 0x03, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06}
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

// 添加更多的测试用例...
