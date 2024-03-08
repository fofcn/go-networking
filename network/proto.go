package network

import (
	"bytes"
	"encoding/binary"
	"errors"
)

type Frame struct {
	Version   uint16
	CmdType   CommandType
	HeaderLen uint16
	Header    interface{}
	Payload   []byte
}

type Codec interface {
	Encode() ([]byte, error)
	Decode(data []byte) (Frame, error)
}

type HeaderCodec interface {
	Encode(header interface{}) ([]byte, error)
	Decode(data []byte) (interface{}, error)
}

var cmdFactory = newCommandFactory()

func AddCodec(cmdType CommandType, headerCodec HeaderCodec) {
	cmdFactory.addCmdCodec(cmdType, headerCodec)
}

func Encode(frame *Frame) ([]byte, error) {
	buf := new(bytes.Buffer)

	encodeVersion(frame, buf)
	encodeCmdType(frame, buf)

	// 编码SubHeader数据
	subHeaderCodec, err := cmdFactory.getCmdCodec(frame.CmdType)
	if err != nil {
		return nil, err
	}

	// 编码SubHeader数据
	subHeaderData, err := subHeaderCodec.Encode(frame.Header)
	if err != nil {
		return nil, err
	}
	frame.HeaderLen = uint16(len(subHeaderData))
	encodeHeaderLen(frame, buf)
	buf.Write(subHeaderData)

	buf.Write(frame.Payload)

	return buf.Bytes(), nil
}

func Decode(frame []byte) (*Frame, error) {
	if len(frame) < 2 {
		return nil, errors.New("frame too short to decode")
	}

	buf := bytes.NewReader(frame)

	proto := &Frame{}

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

	headerCodec, err := cmdFactory.getCmdCodec(proto.CmdType)
	if err != nil {
		return nil, err
	}

	header, err := headerCodec.Decode(varintHeaderData)
	if err != nil {
		return nil, err
	}
	proto.Header = header
	proto.Payload = make([]byte, buf.Len())
	if _, err := buf.Read(proto.Payload); err != nil {
		return nil, errors.New("failed to read payload")
	}

	return proto, nil
}

func encodeVersion(proto *Frame, buf *bytes.Buffer) {
	encodeIntBuf(uint64(proto.Version), buf)
}

func encodeCmdType(proto *Frame, buf *bytes.Buffer) {
	encodeIntBuf(uint64(proto.CmdType), buf)
}

func encodeHeaderLen(proto *Frame, buf *bytes.Buffer) {
	encodeIntBuf(uint64(proto.HeaderLen), buf)
}

func encodeIntBuf(variable uint64, buf *bytes.Buffer) {
	cmdBuf := EncodeInteger(variable)
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

func EncodeInteger(variable uint64) []byte {
	var cmdBuf [binary.MaxVarintLen64]byte
	encodeLen := binary.PutUvarint(cmdBuf[:], variable)
	return cmdBuf[:encodeLen]
}

func DecodeInteger(buf *bytes.Reader) (uint64, error) {
	variable, err := binary.ReadUvarint(buf)
	if err != nil {
		return 0, errors.New("failed to decode command type, invalid bytes")
	}
	return variable, nil
}

type commandFactory struct {
	commandToCodec map[CommandType]HeaderCodec
}

func newCommandFactory() *commandFactory {
	return &commandFactory{
		commandToCodec: make(map[CommandType]HeaderCodec),
	}
}

func (cmdFactory *commandFactory) addCmdCodec(cmdType CommandType, headerCodec HeaderCodec) {
	if _, exists := cmdFactory.commandToCodec[cmdType]; exists {
		// todo logging
	}

	cmdFactory.commandToCodec[cmdType] = headerCodec
}

func (cmdFactory *commandFactory) getCmdCodec(cmdType CommandType) (HeaderCodec, error) {
	var headerCodec HeaderCodec
	var exists bool
	if headerCodec, exists = cmdFactory.commandToCodec[cmdType]; !exists {
		// todo logging
		return nil, errors.New("Codec for command is not find.")
	}

	return headerCodec, nil
}
