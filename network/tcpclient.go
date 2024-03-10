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
	"golang.org/x/sync/semaphore"
)

type TcpClientConfig struct {
	Network string
	Timeout time.Duration
}

type TcpClient struct {
	mux       sync.Mutex
	config    *TcpClientConfig
	connTable map[string]netpoll.Connection
	msgTable  map[uint64]*responseWaiter
}

type responseWaiter struct {
	frame     *Frame
	countdown latch.CountDownLatch
}

func NewTcpClient(config *TcpClientConfig) *TcpClient {
	return &TcpClient{
		config:    config,
		connTable: make(map[string]netpoll.Connection),
		msgTable:  make(map[uint64]*responseWaiter),
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
	// 创建一个超时上下文
	// 设置30秒超时
	// 在完成后释放资源
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	context.WithValue(ctx, "serverAddr", serverAddr)
	conn.SetOnRequest(tcpClient.handleRequest)

	writer := conn.Writer()
	// encode frame
	bytes, err := Encode(frame)
	if err != nil {
		return nil, err
	}

	tcpClient.mux.Lock()
	defer tcpClient.mux.Unlock()
	sem := semaphore.NewWeighted(1)

	respWaiter := &responseWaiter{
		frame:     nil,
		countdown: *latch.NewCountDownLatch(),
	}
	respWaiter.countdown.Add(1)
	tcpClient.msgTable[frame.Sequence] = respWaiter

	// length-value encode
	var lenBytes []byte = make([]byte, binary.MaxVarintLen64)
	encodeLen := binary.PutUvarint(lenBytes, uint64(len(bytes)))
	writer.WriteBinary(lenBytes[:encodeLen])
	cnt, err := writer.WriteBinary(bytes)
	if err != nil || cnt != len(bytes) {
		return nil, err
	}

	err = writer.Flush()
	if err != nil {
		return nil, err
	}
	done := make(chan struct{})
	select {
	case <-done:
		// WaitGroup 完成
		fmt.Println("All goroutines have finished")
	case <-time.After(30 * time.Second):
		// 超时后此分支将执行
		fmt.Println("Timed out waiting for goroutines to finish")
	}
	sem.Acquire(ctx, 1)
	respWaiter.countdown.WaitWithTimeout(30 * time.Second)
	defer respWaiter.countdown.Close()
	if resp, exists := tcpClient.msgTable[frame.Sequence]; exists {
		return resp.frame, nil
	}

	return nil, errors.New("wait response timeout")
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

	conn.SetOnRequest()
	tcpClient.connTable[serverAddr] = conn

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

func (tcpClient *TcpClient) handleRequest(ctx context.Context, conn netpoll.Connection) (err error) {
	reader := conn.Reader()
	len, err := binary.ReadUvarint(reader)
	if err != nil {
		fmt.Printf("%s", err)
		if err == io.EOF {
			conn.Close()
			delete(tcpClient.connTable, ctx.Value(""))
		}
	}
	data, err := reader.ReadBinary(int(len))
	if err != nil {
		fmt.Printf("%s", err)
		if err == io.EOF {
			conn.Close()
			delete(tcpClient.connTable, ctx.Value("serverAddr"))
		}
	}
	frame, err := Decode(data)
	if err != nil {
		fmt.Printf("%s", err)
	}
	if _, exists := tcpClient.msgTable[frame.Sequence]; !exists {
		fmt.Printf("what's wrong? frame sequence not matched with sequence no.: %d", frame.Sequence)
	} else {
		respWaiter := tcpClient.msgTable[frame.Sequence]
		// todo notify waiting client
		respWaiter.countdown.Done()
	}
}
