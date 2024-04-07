package network

import (
	"context"
	"encoding/binary"
	"errors"
	"go-networking/log"
	"io"
	"time"

	"github.com/cloudwego/netpoll"
)

type TcpServerConfig struct {
	Network string
	Addr
}

type TcpServer struct {
	config       *TcpServerConfig
	processors   map[CommandType]Processor
	interceptors []RequestInterceptor
	listener     netpoll.Listener
	eventLoop    netpoll.EventLoop
	pollerNum    int
	connKeyTable map[string]*ConnCtx
	CManager     *ConnManager
}

func (s *TcpServer) AddConnKey(Id string, connKey *ConnCtx) {
	s.connKeyTable[Id] = connKey
}

func NewTcpServer(config *TcpServerConfig) (*TcpServer, error) {
	tcpServer := TcpServer{
		processors:   make(map[CommandType]Processor),
		interceptors: make([]RequestInterceptor, 0),
		config:       config,
		connKeyTable: make(map[string]*ConnCtx),
		CManager:     NewConnManager(),
	}
	return &tcpServer, nil
}

func (s *TcpServer) Init() error {
	log.Info("start tcp server")

	s.config.Network = "tcp"
	s.pollerNum = 2

	netpoll.SetNumLoops(s.pollerNum)

	address := s.config.Addr.Host + ":" + s.config.Port
	listener, err := netpoll.CreateListener(s.config.Network, address)
	if err != nil {
		return err
	}
	s.listener = listener

	eventLoop, err := netpoll.NewEventLoop(
		s.handle,
		netpoll.WithOnPrepare(prepare),
		netpoll.WithOnConnect(connect),
		netpoll.WithReadTimeout(time.Second))
	if err != nil {
		listener.Close()
		return err
	}

	s.eventLoop = eventLoop
	log.Info("started tcp server")
	return nil

}

func (s *TcpServer) Start() error {
	err := s.eventLoop.Serve(s.listener)
	if err != nil {
		return err
	}

	return nil
}

func (s *TcpServer) Stop() error {
	log.Info("TCP Server stop")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	err := s.eventLoop.Shutdown(ctx)
	return err
}

func (s *TcpServer) AddProcessor(cmdType CommandType, process Processor) {
	log.Info("Adding processor")
	s.processors[cmdType] = process
}

func (s *TcpServer) AddInterceptor(requestInterceptor RequestInterceptor) {
	s.interceptors = append(s.interceptors, requestInterceptor)
}

func prepare(connection netpoll.Connection) context.Context {
	return context.Background()
}

func close(connection netpoll.Connection) error {
	log.Infof("[%v] connection closed\n", connection.RemoteAddr())
	return nil
}

func connect(ctx context.Context, connection netpoll.Connection) context.Context {
	log.Infof("[%v] connection established\n", connection.RemoteAddr())
	connection.AddCloseCallback(close)
	return ctx
}

func (s *TcpServer) handle(ctx context.Context, connection netpoll.Connection) error {
	reader, writer := connection.Reader(), connection.Writer()
	readLen, err := binary.ReadUvarint(reader)
	if err != nil {
		log.Errorf("%s", err)
		if err == io.EOF {
			defer reader.Release()
		}

		return err
	}

	data, err := reader.ReadBinary(int(readLen))
	if err != nil {
		log.Errorf("%s", err)
		if err == io.EOF {
			defer reader.Release()
		}
	}

	req, err := Decode(LVBasedCodec, data)
	if err != nil {
		return err
	}

	log.Infof("server recv frame sequence: %d", req.Seq)

	if len(s.interceptors) != 0 {
		for _, interceptor := range s.interceptors {
			// todo add client address
			interceptor.OnRequest("", req)
		}
	}

	if processor, ok := s.processors[req.CmdType]; ok {
		resp, err := processor.Process(&Conn{
			Connection: connection,
		}, req)
		if err != nil {
			return err
		}
		if len(s.interceptors) != 0 {
			for _, interceptor := range s.interceptors {
				// todo add client direction
				interceptor.OnResponse("", req, resp)
			}
		}

		respData, err := Encode(LVBasedCodec, resp)
		if err != nil {
			return err
		}

		_, err = writer.WriteBinary(respData)
		if err != nil {
			return err
		}

		err = writer.Flush()
		if err != nil {
			return err
		}
	} else {
		return errors.New("command processor cannot be found")
	}

	return nil
}
