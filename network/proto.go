package network

import "go-networking/network"

type VersionType uint16

const (
	VERSION_1 VersionType = iota + 1
)

type Frame struct {
	Version VersionType
	CmdType CommandType
	Seq     uint64
	HLen    uint16
	Header  interface{}
	Payload []byte
}

func NewFrame(cmdType CommandType, h interface{}, payload []byte) *Frame {
	return &Frame{
		Version: network.VERSION_1,
		CmdType: cmdType,
		Header:  h,
		Payload: payload,
	}
}

const (
	LVBasedCodec = iota
)

func Encode(codecChoice int, frame *Frame) ([]byte, error) {
	if codecChoice == LVBasedCodec {
		codec := NewLVCodec()
		data, err := codec.Encode(frame)
		if err != nil {
			return nil, err
		}

		return data, nil
	} else {
		return nil, nil
	}
}

func Decode(codecChoice int, data []byte) (*Frame, error) {
	if codecChoice == LVBasedCodec {
		codec := NewLVCodec()
		frame, err := codec.Decode(data)
		if err != nil {
			return nil, err
		}

		return frame, nil
	} else {
		return nil, nil
	}
}
