package network

type Frame struct {
	Version   uint16
	CmdType   CommandType
	Sequence  uint64
	HeaderLen uint16
	Header    interface{}
	Payload   []byte
}

const (
	LengthValueBasedCodec = iota
)

func Encode(codecChoice int, frame *Frame) ([]byte, error) {
	if codecChoice == LengthValueBasedCodec {
		codec := NewLengthValueCodec()
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
	if codecChoice == LengthValueBasedCodec {
		codec := NewLengthValueCodec()
		frame, err := codec.Decode(data)
		if err != nil {
			return nil, err
		}

		return frame, nil
	} else {
		return nil, nil
	}
}
