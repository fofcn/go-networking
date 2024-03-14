package network

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
)

type Codec interface {
	Encode(frame *Frame) ([]byte, error)
	Decode(data []byte) (*Frame, error)
	GetLength() int
}

type HeaderCodec interface {
	Encode(header interface{}) ([]byte, error)
	Decode(data []byte) (interface{}, error)
}

type LVCodec struct {
}

func NewLVCodec() *LVCodec {
	return &LVCodec{}
}

func AddHeaderCodec(cmdType CommandType, headerCodec HeaderCodec) {
	cmdFactory.addCmdCodec(cmdType, headerCodec)
}

func (codec *LVCodec) Encode(frame *Frame) ([]byte, error) {
	buf := new(bytes.Buffer)
	encodeVersion(frame, buf)
	encodeCmdType(frame, buf)
	encodeSeq(frame, buf)

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
	frame.HLen = uint16(len(subHeaderData))
	encodeHLen(frame, buf)
	buf.Write(subHeaderData)

	buf.Write(frame.Payload)
	var lengthBytes []byte = make([]byte, binary.MaxVarintLen32)
	encodeLen := binary.PutUvarint(lengthBytes, uint64(buf.Len()))

	finalBuf := new(bytes.Buffer)
	finalBuf.Write(lengthBytes[:encodeLen])
	finalBuf.Write(buf.Bytes())
	return finalBuf.Bytes(), nil
}

func (codec *LVCodec) Decode(data []byte) (*Frame, error) {
	if len(data) < 2 {
		return nil, errors.New("frame too short to decode")
	}

	buf := bytes.NewReader(data)

	frame := &Frame{}

	if err := decodeVersion(buf, &frame.Version); err != nil {
		return nil, err
	}

	if err := decodeCmdType(buf, &frame.CmdType); err != nil {
		return nil, err
	}

	if err := decodeSeq(buf, &frame.Seq); err != nil {
		return nil, err
	}

	if err := decodeHLen(buf, &frame.HLen); err != nil {
		return nil, err
	}

	varintHeaderData := make([]byte, frame.HLen)
	if n, err := buf.Read(varintHeaderData); err != nil || n != int(frame.HLen) {
		return nil, errors.New("failed to read the correct subheader length")
	}

	headerCodec, err := cmdFactory.getCmdCodec(frame.CmdType)
	if err != nil {
		return nil, err
	}

	header, err := headerCodec.Decode(varintHeaderData)
	if err != nil {
		return nil, err
	}
	frame.Header = header
	frame.Payload = make([]byte, buf.Len())
	if _, err := buf.Read(frame.Payload); err != nil {
		return nil, errors.New("failed to read payload")
	}

	return frame, nil
}

func encodeVersion(frame *Frame, buf *bytes.Buffer) {
	encodeIntBuf(uint64(frame.Version), buf)
}

func encodeCmdType(frame *Frame, buf *bytes.Buffer) {
	encodeIntBuf(uint64(frame.CmdType), buf)
}

func encodeSeq(frame *Frame, buf *bytes.Buffer) {
	encodeIntBuf(frame.Seq, buf)
}

func encodeHLen(frame *Frame, buf *bytes.Buffer) {
	encodeIntBuf(uint64(frame.HLen), buf)
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

func decodeSeq(buf *bytes.Reader, sequence *uint64) error {
	decSeq, err := binary.ReadUvarint(buf)
	if err != nil {
		return errors.New("failed to decode command type, invalid bytes")
	}

	*sequence = decSeq
	return nil
}

func decodeHLen(buf *bytes.Reader, headerLen *uint16) error {
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

var cmdFactory = newCommandFactory()

type commandFactory struct {
	commandToCodec map[CommandType]HeaderCodec
}

func newCommandFactory() *commandFactory {
	return &commandFactory{
		commandToCodec: make(map[CommandType]HeaderCodec),
	}
}

func (cmdFactory *commandFactory) addCmdCodec(cmdType CommandType, headerCodec HeaderCodec) {
	cmdFactory.commandToCodec[cmdType] = headerCodec
}

func (cmdFactory *commandFactory) getCmdCodec(cmdType CommandType) (HeaderCodec, error) {
	var headerCodec HeaderCodec
	var exists bool
	if headerCodec, exists = cmdFactory.commandToCodec[cmdType]; !exists {
		return nil, fmt.Errorf("codec for command is not find, cmd type: %d", cmdType)
	}

	return headerCodec, nil
}
