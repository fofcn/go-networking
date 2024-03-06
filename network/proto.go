package network

import (
	"bytes"
	"encoding/binary"
	"errors"
)

type Proto struct {
	Version   uint16
	CmdType   CommandType
	HeaderLen uint16
	Header    []byte
	Payload   []byte
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

	encodeVersion(proto, buf)
	encodeCmdType(proto, buf)

	// 编码SubHeader数据
	subHeaderCodec, err := cmdFactory.getCmdCodec(proto.CmdType)
	if err != nil {
		return nil, err
	}

	// 编码SubHeader数据
	subHeaderData, err := subHeaderCodec.Encode()
	if err != nil {
		return nil, err
	}
	proto.HeaderLen = uint16(len(subHeaderData))
	encodeHeaderLen(proto, buf)
	buf.Write(subHeaderData)

	buf.Write(proto.Payload)

	return buf.Bytes(), nil
}

func Decode(frame []byte) (*Proto, error) {
	if len(frame) < 2 {
		return nil, errors.New("frame too short to decode")
	}

	buf := bytes.NewReader(frame)

	proto := &Proto{}

	if err := decodeVersion(buf, &proto.Version); err != nil {
		return nil, err
	}

	if err := decodeCmdType(buf, &proto.CmdType); err != nil {
		return nil, err
	}

	if err := decodeHeaderLen(buf, &proto.HeaderLen); err != nil {
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

	if err := varintHeader.Decode(varintHeaderData); err != nil {
		return nil, err
	}

	proto.Payload = make([]byte, buf.Len())
	if _, err := buf.Read(proto.Payload); err != nil {
		return nil, errors.New("failed to read payload")
	}

	return proto, nil
}

func encodeVersion(proto *Proto, buf *bytes.Buffer) {
	encodeIntBuf(uint64(proto.Version), buf)
}

func encodeCmdType(proto *Proto, buf *bytes.Buffer) {
	encodeIntBuf(uint64(proto.CmdType), buf)
}

func encodeHeaderLen(proto *Proto, buf *bytes.Buffer) {
	encodeIntBuf(uint64(proto.HeaderLen), buf)
}

func encodeIntBuf(variable uint64, buf *bytes.Buffer) {
	cmdBuf := encodeInteger(variable)
	buf.Write(cmdBuf)
}

func decodeVersion(buf *bytes.Reader, version *uint16) error {
	decVersion, err := binary.ReadUvarint(buf)
	if err != nil {
		return errors.New("failed to decode version, invalid bytes")
	}

	*version = uint16(decVersion)
	return nil
}

func decodeCmdType(buf *bytes.Reader, version *CommandType) error {
	decVersion, err := binary.ReadUvarint(buf)
	if err != nil {
		return errors.New("failed to decode command type, invalid bytes")
	}

	*version = CommandType(decVersion)
	return nil
}

func decodeHeaderLen(buf *bytes.Reader, headerLen *uint16) error {
	decVersion, err := binary.ReadUvarint(buf)
	if err != nil {
		return errors.New("failed to decode command type, invalid bytes")
	}

	*headerLen = uint16(decVersion)
	return nil
}

func encodeInteger(variable uint64) []byte {
	var cmdBuf [binary.MaxVarintLen64]byte
	encodeLen := binary.PutUvarint(cmdBuf[:], variable)
	return cmdBuf[:encodeLen]
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
	if codec, exists = cmdFactory.commandToCodec[cmdType]; !exists {
		// todo logging
		return nil, errors.New("Codec for command is not find.")
	}

	return codec, nil
}
