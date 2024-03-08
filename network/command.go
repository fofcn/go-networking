package network

type CommandType int32

const (
	CONN CommandType = iota
	CONNACK
)

type ConnCmd struct {
	ConnBase
}

type ConnAckCmd struct {
	ConnBase
}

type ConnBase struct {
	KeyLen uint32
	Key    string
}
