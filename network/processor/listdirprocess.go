package processor

import (
	"go-networking/network"
)

type ListdireProcessor struct {
	TcpSrv *network.TcpServer
}

func (lp *ListdireProcessor) Process(conn *network.Conn, frame *network.Frame) (*network.Frame, error) {
	return nil, nil
}
