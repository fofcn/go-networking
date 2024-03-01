package network

type Packet struct {
}

type Addr struct {
	Host string
	Port string
}

type Conn interface {
}

type Processor interface {
	Process(conn *Conn, packet *Packet)
}

type Lifecycle interface {
	Init() error
	Start() error
	Stop() error
}

type Server interface {
	Lifecycle
	AddProcessor(cmdType CommandType, process Processor)
}

type Client interface {
	Lifecycle
	Send(addr *Addr, packet *Packet)
}
