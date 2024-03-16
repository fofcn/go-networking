package processor

import (
	"go-networking/network"
	"time"
)

type PingProcessor struct {
	TcpServer *network.TcpServer
}

func (pp *PingProcessor) Process(conn *network.Conn, frame *network.Frame) (*network.Frame, error) {
	header := frame.Header.(*network.PingHeader)
	pp.TcpServer.CManager.Ping(header.Id, header.Timestamp)

	return network.NewFrame(network.PONG,
		&network.PongHeader{
			Timestamp: time.Now().Unix(),
		},
		nil), nil
}
