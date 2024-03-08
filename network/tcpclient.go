package network

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
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
	msgTable  map[uint64]*Frame
}

func NewTcpClient(config *TcpClientConfig) *TcpClient {
	return &TcpClient{
		config:    config,
		connTable: make(map[string]netpoll.Connection),
		msgTable:  make(map[uint64]*Frame),
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

func (tcpClient *TcpClient) SendSync(serverAddr string, frame *Frame) (*Frame, error) {
	conn, err := tcpClient.getOrCreateConnection(tcpClient.config.Network, serverAddr, tcpClient.config.Timeout)
	if err != nil {
		return nil, err
	}

	writer := conn.Writer()
	// encode frame
	bytes, err := Encode(frame)
	if err != nil {
		return nil, err
	}

	cnt, err := writer.WriteBinary(bytes)
	if err != nil || cnt != len(bytes) {
		return nil, err
	}

	err = writer.Flush()
	if err != nil {
		return nil, err
	}

	tcpClient.mux.Lock()
	defer tcpClient.mux.Unlock()
	tcpClient.msgTable[frame.Sequence] = nil

	return nil, nil
}

func (tcpClient *TcpClient) SendAsync(serverAddr string, packet *Frame) error {
	return nil
}

func (tcpClient *TcpClient) SendOnce(serverAddr string, packet *Frame) error {
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
	if exists && conn.IsActive() {
		return conn, nil
	}

	if !conn.IsActive() {
		err := conn.Close()
		fmt.Printf("connection was not active, close also occured error, please check the error: %s", err)
		delete(tcpClient.connTable, serverAddr)
	}

	conn, err := tcpClient.createConnection(network, serverAddr, timeout)
	if err != nil {
		return nil, err
	}

	if !conn.IsActive() {
		return nil, errors.New("connection have not actived after created")
	}

	tcpClient.connTable[serverAddr] = conn
	go func() {
		reader := conn.Reader()
		len, err := binary.ReadUvarint(reader)
		if err != nil {
			fmt.Printf("%s", err)
			if err == io.EOF {
				conn.Close()
				delete(tcpClient.connTable, serverAddr)
			}
		}
		data, err := reader.ReadBinary(int(len))
		if err != nil {
			fmt.Printf("%s", err)
			if err == io.EOF {
				conn.Close()
				delete(tcpClient.connTable, serverAddr)
			}
		}
		frame, err := Decode(data)
		if err != nil {
			fmt.Printf("%s", err)
		}
		if _, exists := tcpClient.msgTable[frame.Sequence]; !exists {
			fmt.Printf("what's wrong? frame sequence not matched with sequence no.: %d", frame.Sequence)
		} else {
			tcpClient.msgTable[frame.Sequence] = frame
			// todo notify waiting client
		}
	}()

	return conn, nil
}

func (TcpClient *TcpClient) createConnection(network string, serverAddr string, timeout time.Duration) (netpoll.Connection, error) {
	conn, err := netpoll.DialConnection(network, serverAddr, timeout)
	if err != nil {
		return nil, err
	}
	conn.AddCloseCallback(func(connection netpoll.Connection) error {
		fmt.Printf("[%v] connection closed\n", connection.RemoteAddr())
		return nil
	})

	return conn, err
}
