package network

import (
	"context"
	"fmt"
	"runtime"
	"time"

	"github.com/cloudwego/netpoll"
)

type TcpServerConfig struct {
	Network string
	Addr
}

type TcpServer struct {
	config     *TcpServerConfig
	processors map[CommandType]*Processor
	listener   netpoll.Listener
	eventLoop  netpoll.EventLoop
	pollerNum  int
}

func NewTcpServer(config *TcpServerConfig) (*TcpServer, error) {
	tcpServer := TcpServer{
		processors: make(map[CommandType]*Processor),
		config:     config,
	}
	return &tcpServer, nil
}

func (tcpServer *TcpServer) Init() error {
	fmt.Println("start tcp server")

	tcpServer.config.Network = "tcp"
	tcpServer.pollerNum = 2

	// netpoll.SetNumLoops(tcpServer.pollerNum)
	runtime.GOMAXPROCS(tcpServer.pollerNum)
	address := tcpServer.config.Addr.Host + ":" + tcpServer.config.Port
	listener, err := netpoll.CreateListener(tcpServer.config.Network, address)
	if err != nil {
		return err
	}
	tcpServer.listener = listener

	eventLoop, err := netpoll.NewEventLoop(
		handle,
		netpoll.WithOnPrepare(prepare),
		netpoll.WithOnConnect(connect),
		netpoll.WithReadTimeout(time.Second))
	if err != nil {
		listener.Close()
		return err
	}

	tcpServer.eventLoop = eventLoop
	fmt.Println("started tcp server")
	return nil

}

func (tcpServer *TcpServer) Start() error {
	err := tcpServer.eventLoop.Serve(tcpServer.listener)
	if err != nil {
		return err
	}

	return nil
}

func (tcpServer *TcpServer) Stop() error {
	fmt.Println("TCP Server stop")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := tcpServer.eventLoop.Shutdown(ctx)
	return err
}

func (tcpServer *TcpServer) AddProcessor(cmdType CommandType, process Processor) {
	fmt.Println("Adding processor")
	tcpServer.processors[cmdType] = &process
}

func (tcpServer *TcpServer) AddInterceptor(requestInterceptor RequestInterceptor) {

}

type ConnProcessor struct {
}

func (ConnProcessor ConnProcessor) Process(conn *Conn, frame *Frame) {

}

func prepare(connection netpoll.Connection) context.Context {
	return context.Background()
}

func close(connection netpoll.Connection) error {
	fmt.Printf("[%v] connection closed\n", connection.RemoteAddr())
	return nil
}

func connect(ctx context.Context, connection netpoll.Connection) context.Context {
	fmt.Printf("[%v] connection established\n", connection.RemoteAddr())
	connection.AddCloseCallback(close)
	return ctx
}

func handle(ctx context.Context, connection netpoll.Connection) error {
	reader, writer := connection.Reader(), connection.Writer()
	defer reader.Release()

	msg, _ := reader.ReadString(reader.Len())
	fmt.Printf("[recv msg] %v\n", msg)

	writer.WriteString(msg)
	writer.Flush()

	return nil
}
