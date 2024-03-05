package network

import (
	"bytes"
	"encoding/binary"
	"errors"
)

type Proto struct {
	Version      uint16
	CmdType      CommandType
	HeaderLen    uint16
	VarintHeader VarintHeader
	Payload      []byte
}

type VarintHeader interface {
	Encode() ([]byte, error)
	Decode([]byte) error
}

var cmdFactory = newCommandFactory()

func AddCodec(cmdType CommandType, varintHeader VarintHeader) {
	cmdFactory.addCmdCodec(cmdType, varintHeader)
}

func Encode(proto *Proto) ([]byte, error) {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, proto.Version)
	binary.Write(buf, binary.BigEndian, proto.CmdType)

	// 编码SubHeader数据
	subHeaderData, err := proto.VarintHeader.Encode()
	if err != nil {
		return nil, err
	}

	proto.HeaderLen = uint16(len(subHeaderData))
	binary.Write(buf, binary.BigEndian, proto.HeaderLen)
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

	if err := binary.Read(buf, binary.BigEndian, &proto.HeaderLen); err != nil {
		return nil, err
	}

	varintHeaderData := make([]byte, proto.HeaderLen)
	if n, err := buf.Read(varintHeaderData); err != nil || n != int(proto.HeaderLen) {
		return nil, errors.New("failed to read the correct subheader length")
	}

	varintHeader, err := cmdFactory.getCmdCodec(proto.CmdType)
	if err != nil {
		return nil, err
	}
	proto.VarintHeader = varintHeader
	if err := proto.VarintHeader.Decode(varintHeaderData); err != nil {
		return nil, err
	}

	proto.Payload = make([]byte, buf.Len())
	if _, err := buf.Read(proto.Payload); err != nil {
		return nil, errors.New("failed to read payload")
	}

	return proto, nil
}

type commandFactory struct {
	commandToCodec map[CommandType]VarintHeader
}

func newCommandFactory() *commandFactory {
	return &commandFactory{
		commandToCodec: make(map[CommandType]VarintHeader),
	}
}

func (cmdFactory *commandFactory) addCmdCodec(cmdType CommandType, header VarintHeader) {
	if _, exists := cmdFactory.commandToCodec[cmdType]; exists {
		// todo logging
	}

	cmdFactory.commandToCodec[cmdType] = header
}

func (cmdFactory *commandFactory) getCmdCodec(cmdType CommandType) (VarintHeader, error) {
	var codec VarintHeader
	var exists bool
	if codec, exists = cmdFactory.commandToCodec[cmdType]; exists {
		// todo logging
		return nil, errors.New("Codec for command is not find.")
	}

	return codec, nil
}
