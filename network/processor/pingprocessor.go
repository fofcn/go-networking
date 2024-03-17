package processor

import (
	"go-networking/network"
	"go-networking/network/codec"
	"time"
)

type PingProcessor struct {
	TcpServer *network.TcpServer
}

func (pp *PingProcessor) Process(conn *network.Conn, frame *network.Frame) (*network.Frame, error) {
	header := frame.Header.(*codec.PingHeader)
	pp.TcpServer.CManager.Ping(header.Id, header.Timestamp)

	return network.NewFrame(network.PONG,
		&codec.PongHeader{
			Timestamp: time.Now().Unix(),
		},
		nil), nil
}
