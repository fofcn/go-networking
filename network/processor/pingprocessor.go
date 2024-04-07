package processor

import (
	"errors"
	"go-networking/log"
	"go-networking/network"
	"go-networking/network/codec"
	"time"
)

type PingProcessor struct {
	tcpSrv *network.TcpServer
}

func NewPingProcs(tcpSrv *network.TcpServer) *PingProcessor {
	return &PingProcessor{
		tcpSrv: tcpSrv,
	}
}

func (pp *PingProcessor) Process(conn *network.Conn, frame *network.Frame) (*network.Frame, error) {
	header := frame.Header.(*codec.PingHeader)
	err := pp.tcpSrv.CManager.Ping(header.Id, header.Timestamp)
	if err != nil && errors.Is(err, codec.Invalid_Ping_Frame) {
		log.Info("ignore this ping frame")
		return nil, err
	} else if err != nil {
		pp.tcpSrv.CManager.Delete(header.Id)
		return nil, err
	}

	return network.NewFrame(network.PONG,
		&codec.PongHeader{
			Timestamp: time.Now().Unix(),
		},
		nil), nil
}
