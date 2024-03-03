package network

import (
	"fmt"
	"sync"
	"time"

	"github.com/cloudwego/netpoll"
)

type TcpClientConfig struct {
	Network string
	Timeout time.Duration
}

type TcpClient struct {
	mux       sync.Mutex
	config    *TcpClientConfig
	connTable map[string]netpoll.Connection
}

func NewTcpClient(config *TcpClientConfig) *TcpClient {
	return &TcpClient{
		config:    config,
		connTable: make(map[string]netpoll.Connection),
	}
}

func (tcpClient *TcpClient) Init() error {

	return nil
}

func (tcpClient *TcpClient) Start() error {
	return nil
}

func (tcpClient *TcpClient) Stop() error {
	return nil
}

func (tcpClient *TcpClient) SendSync(serverAddr string, packet *Packet) (*Packet, error) {
	conn, err := tcpClient.getOrCreateConnection(tcpClient.config.Network, serverAddr, tcpClient.config.Timeout)
	if err != nil {
		return nil, err
	}

	writer := conn.Writer()
	// encode Packet

	var packetBuf []byte
	cnt, err := writer.WriteBinary()
	if err != nil || cnt != len(packetBuf) {
		return nil, err
	}

	err = writer.Flush()
	if err != nil {
		return nil, err
	}

	reader := conn.Reader()
	reader.ReadBinary()

	return nil, nil
}

func (tcpClient *TcpClient) SendAsync(serverAddr string, packet *Packet) error {
	return nil
}

func (tcpClient *TcpClient) SendOnce(serverAddr string, packet *Packet) error {
	return nil
}

func (tcpClient *TcpClient) AddProcessor(commandType CommandType, processor Processor) {

}

func (tcpClient *TcpClient) AddInterceptor(requestInterceptor RequestInterceptor) {

}

func (tcpClient *TcpClient) getOrCreateConnection(network string, serverAddr string, timeout time.Duration) (netpoll.Connection, error) {
	tcpClient.mux.Lock()
	defer tcpClient.mux.Unlock()
	conn, exists := tcpClient.connTable[serverAddr]
	if exists {
		return conn, nil
	}

	conn, err := tcpClient.createConnection(network, serverAddr, timeout)
	if err != nil {
		return nil, err
	}

	tcpClient.connTable[serverAddr] = conn
	return conn, nil
}

func (TcpClient *TcpClient) createConnection(network string, serverAddr string, timeout time.Duration) (netpoll.Connection, error) {

	conn, err := netpoll.DialConnection(network, serverAddr, timeout)
	if err != nil {
		return nil, err
	}
	conn.AddCloseCallback(func(connection netpoll.Connection) (netpoll.Connection, error) {
		fmt.Printf("[%v] connection closed\n", connection.RemoteAddr())
		return nil, nil
	})

	return conn, err
}
