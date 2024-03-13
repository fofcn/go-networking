package network

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/cloudwego/netpoll"
	"github.com/quintans/toolkit/latch"
)

type TcpClientConfig struct {
	Network string
	Timeout time.Duration
}

type TcpClient struct {
	mux       sync.Mutex
	config    *TcpClientConfig
	connTable map[string]netpoll.Connection
	msgTable  map[uint64]ResponseFuture
}

type responseWaiter struct {
	frame     *Frame
	countdown latch.CountDownLatch
}

func NewTcpClient(config *TcpClientConfig) *TcpClient {
	return &TcpClient{
		config:    config,
		connTable: make(map[string]netpoll.Connection),
		msgTable:  make(map[uint64]ResponseFuture),
	}
}

func (tcpClient *TcpClient) Init() error {

	return nil
}

func (tcpClient *TcpClient) Start() error {
	return nil
}

func (tcpClient *TcpClient) Stop() error {
	for _, conn := range tcpClient.connTable {
		if conn.IsActive() {
			conn.Close()
		}
	}

	for _, future := range tcpClient.msgTable {
		future.Close()
	}
	return nil
}

type contextKey string

func (tcpClient *TcpClient) SendSync(serverAddr string, frame *Frame, timeout time.Duration) (*Frame, error) {
	conn, err := tcpClient.getOrCreateConnection(tcpClient.config.Network, serverAddr, tcpClient.config.Timeout)
	if err != nil {
		return nil, err
	}
	// 创建一个超时上下文
	// 设置30秒超时
	// 在完成后释放资源
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	var key contextKey = "serverAddr"
	ctx = context.WithValue(ctx, key, serverAddr)
	conn.SetOnRequest(tcpClient.handleRequest)

	writer := conn.Writer()
	// encode frame
	bytes, err := Encode(LengthValueBasedCodec, frame)
	if err != nil {
		return nil, err
	}

	rf := NewResponseFuture(frame.Sequence, timeout)
	defer rf.Close()
	tcpClient.addSeqFuture(frame.Sequence, rf)

	cnt, err := writer.WriteBinary(bytes)
	if err != nil || cnt != len(bytes) {
		return nil, errors.New("send failed")
	}

	err = writer.Flush()
	if err != nil {
		return nil, err
	}

	respFrame, err := rf.Wait()
	if err != nil {
		return nil, err
	}
	delete(tcpClient.msgTable, frame.Sequence)
	return respFrame, nil
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
	if exists {
		if conn.IsActive() {
			return conn, nil
		} else {
			err := conn.Close()
			fmt.Printf("connection was not active, close also occured error, please check the error: %s", err)
			delete(tcpClient.connTable, serverAddr)
		}
	}

	conn, err := tcpClient.createConnection(network, serverAddr, timeout)
	if err != nil {
		return nil, err
	}

	if !conn.IsActive() {
		return nil, errors.New("connection have not actived after created")
	}

	tcpClient.connTable[serverAddr] = conn

	return conn, nil
}

func (TcpClient *TcpClient) createConnection(network string, serverAddr string, timeout time.Duration) (netpoll.Connection, error) {
	conn, err := netpoll.DialConnection(network, serverAddr, timeout)
	if err != nil {
		return nil, err
	}
	conn.AddCloseCallback(TcpClient.closeConnectionCallback)
	return conn, err
}

func (tcpClient *TcpClient) handleRequest(ctx context.Context, conn netpoll.Connection) (err error) {
	reader := conn.Reader()
	len, err := binary.ReadUvarint(reader)
	if err != nil {
		fmt.Printf("%s", err)
		if err == io.EOF {
			conn.Close()
			serevrAddr := ctx.Value("serverAddr")
			if addr, ok := serevrAddr.(string); ok {
				delete(tcpClient.connTable, addr)
			}
		}

		return err
	}
	data, err := reader.ReadBinary(int(len))
	if err != nil {
		fmt.Printf("%s", err)
		if err == io.EOF {
			conn.Close()
			serevrAddr := ctx.Value("serverAddr")
			if addr, ok := serevrAddr.(string); ok {
				delete(tcpClient.connTable, addr)
			}
		}
	}
	frame, err := Decode(LengthValueBasedCodec, data)
	if err != nil {
		fmt.Printf("%s", err)
		return err
	}
	fmt.Printf("received frame sequence no.: %d", frame.Sequence)
	if _, exists := tcpClient.msgTable[frame.Sequence]; !exists {
		fmt.Printf("what's wrong? frame sequence not matched with sequence no.: %d", frame.Sequence)
	} else {
		rf := tcpClient.msgTable[frame.Sequence]
		rf.Add(frame)
	}

	return nil
}

func (tcpClient *TcpClient) closeConnectionCallback(conn netpoll.Connection) error {
	fmt.Printf("[%v] connection closed\n", conn.RemoteAddr())
	addr := conn.RemoteAddr()
	conn.Close()
	delete(tcpClient.connTable, addr.String())
	return nil
}

func (tcpClient *TcpClient) addSeqFuture(seq uint64, rf ResponseFuture) {
	tcpClient.mux.Lock()
	defer tcpClient.mux.Unlock()

	tcpClient.msgTable[seq] = rf
}
