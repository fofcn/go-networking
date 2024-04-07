package processor

import (
	"go-networking/network"
)

type ListdireProcessor struct {
	tcpSrv *network.TcpServer
}

func NewListdireProcs(tcpSrv *network.TcpServer) *ListdireProcessor {
	return &ListdireProcessor{
		tcpSrv: tcpSrv,
	}
}

func (lp *ListdireProcessor) Process(conn *network.Conn, frame *network.Frame) (*network.Frame, error) {
	return nil, nil
}
