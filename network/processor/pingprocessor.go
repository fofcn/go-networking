package processor

import "go-networking/network"

type PingProcessor struct {
	TcpServer *network.TcpServer
}

func (pp *PingProcessor) Process(conn *network.Conn, frame *network.Frame) (*network.Frame, error) {
	frame.
}
