package network

import (
	"bytes"
	"encoding/binary"
	"errors"
)

type Proto struct {
	Version      uint16
	CmdType      CommandType
	SubHeaderLen uint16
	VarintHeader VarintHeader
	Payload      []byte
}

type VarintHeader interface {
	Encode() ([]byte, error)
	Decode([]byte) error
}

func Encode(proto *Proto) ([]byte, error) {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, proto.Version)
	binary.Write(buf, binary.BigEndian, proto.CmdType)

	// 编码SubHeader数据
	subHeaderData, err := VarintHeader.Encode()
	if err != nil {
		return nil, err
	}

	proto.SubHeaderLen = uint16(len(subHeaderData))
	binary.Write(buf, binary.BigEndian, proto.SubHeaderLen)
	buf.Write(subHeaderData)

	buf.Write(proto.Payload)

	return buf.Bytes(), nil
}

func Decode(frame []byte) (*Proto, error) {
	if len(frame) < 5 {
		return nil, errors.New("frame too short to decode")
	}

	buf := bytes.NewReader(frame)

	proto := &Proto{}

	if err := binary.Read(buf, binary.BigEndian, &proto.Version); err != nil {
		return nil, err
	}

	if err := binary.Read(buf, binary.BigEndian, &proto.CmdType); err != nil {
		return nil, err
	}

	if err := binary.Read(buf, binary.BigEndian, &proto.SubHeaderLen); err != nil {
		return nil, err
	}

	subHeaderData := make([]byte, proto.SubHeaderLen)
	if n, err := buf.Read(subHeaderData); err != nil || n != int(proto.SubHeaderLen) {
		return nil, errors.New("failed to read the correct subheader length")
	}

	// 根据CmdType创建对应的SubHeader实例并解码
	switch proto.CmdType {
	case CommandA:
		proto.SubHeader = &SubHeaderTypeA{}
	case CommandB:
		proto.SubHeader = &SubHeaderTypeB{}
	// ... 更多的case ...
	default:
		return nil, errors.New("unknown command type")
	}

	if err := VarintHeader.Decode(subHeaderData); err != nil {
		return nil, err
	}

	proto.Payload = make([]byte, buf.Len())
	if _, err := buf.Read(proto.Payload); err != nil {
		return nil, errors.New("failed to read payload")
	}

	return proto, nil
}
