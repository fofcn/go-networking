package network

import "fmt"

type CommandType int

const (
	CONN CommandType = iota
)

type Packet struct {
}

type Addr struct {
}

type Conn interface {
}

type Processor interface {
	Process(conn *Conn, packet *Packet)
}

type Lifecycle interface {
	Start()
	Stop()
}

type Server interface {
	Lifecycle
	AddProcessor(cmdType CommandType, process Processor)
}

type Client interface {
	Lifecycle
	Send(addr *Addr, packet *Packet)
}

type TcpServer struct {
	processors map[CommandType]*Processor
}

func NewTcpServer() (*TcpServer, error) {
	tcpServer := TcpServer{
		processors: make(map[CommandType]*Processor),
	}
	return &tcpServer, nil
}

func (tcpServer *TcpServer) Start() {
	fmt.Println("TCP Server start")
}

func (tcpServer *TcpServer) Stop() {
	fmt.Println("TCP Server stop")
}

func (tcpServer *TcpServer) AddProcessor(cmdType CommandType, process Processor) {
	fmt.Println("Adding processor")
	tcpServer.processors[cmdType] = &process
}

type ConnProcessor struct {
}

func (ConnProcessor ConnProcessor) Process(conn *Conn, packet *Packet) {

}
