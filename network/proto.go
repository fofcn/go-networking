package network

type Frame struct {
	Version uint16
	CmdType CommandType
	Seq     uint64
	HLen    uint16
	Header  interface{}
	Payload []byte
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
