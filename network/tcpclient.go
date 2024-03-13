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
)

type TcpClientConfig struct {
	Network string
	Timeout time.Duration
}

type TcpClient struct {
	mux       sync.Mutex
	config    *TcpClientConfig
	connTable map[string]netpoll.Connection
	respTable map[uint64]ResponseFuture
	ticker    *time.Ticker
}

func NewTcpClient(config *TcpClientConfig) *TcpClient {
	return &TcpClient{
		config:    config,
		connTable: make(map[string]netpoll.Connection),
		respTable: make(map[uint64]ResponseFuture),
	}
}

func (c *TcpClient) Init() error {

	return nil
}

func (c *TcpClient) Start() error {
	c.cleanupResponseFutures()
	return nil
}

func (c *TcpClient) Stop() error {
	for _, conn := range c.connTable {
		if conn.IsActive() {
			conn.Close()
		}
	}

	for _, future := range c.respTable {
		future.Close()
	}

	c.ticker.Stop()
	return nil
}

type contextKey string

func (c *TcpClient) SendSync(serverAddr string, frame *Frame, timeout time.Duration) (*Frame, error) {
	conn, err := c.getOrCreateConnection(c.config.Network, serverAddr, c.config.Timeout)
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
	conn.SetOnRequest(c.handleRequest)

	writer := conn.Writer()
	// encode frame
	bytes, err := Encode(LengthValueBasedCodec, frame)
	if err != nil {
		return nil, err
	}

	rf := NewResponseFuture(frame.Sequence, timeout)
	defer rf.Close()
	c.addSeqFuture(frame.Sequence, rf)

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
	delete(c.respTable, frame.Sequence)
	return respFrame, nil
}

func (c *TcpClient) SendAsync(serverAddr string, packet *Frame) error {
	return nil
}

func (c *TcpClient) SendOnce(serverAddr string, packet *Frame) error {
	return nil
}

func (c *TcpClient) AddProcessor(commandType CommandType, processor Processor) {

}

func (c *TcpClient) AddInterceptor(requestInterceptor RequestInterceptor) {

}

func (c *TcpClient) getOrCreateConnection(network string, serverAddr string, timeout time.Duration) (netpoll.Connection, error) {
	c.mux.Lock()
	defer c.mux.Unlock()
	conn, exists := c.connTable[serverAddr]
	if exists {
		if conn.IsActive() {
			return conn, nil
		} else {
			err := conn.Close()
			fmt.Printf("connection was not active, close also occured error, please check the error: %s", err)
			delete(c.connTable, serverAddr)
		}
	}

	conn, err := c.createConnection(network, serverAddr, timeout)
	if err != nil {
		return nil, err
	}

	if !conn.IsActive() {
		return nil, errors.New("connection have not actived after created")
	}

	c.connTable[serverAddr] = conn

	return conn, nil
}

func (c *TcpClient) createConnection(network string, serverAddr string, timeout time.Duration) (netpoll.Connection, error) {
	conn, err := netpoll.DialConnection(network, serverAddr, timeout)
	if err != nil {
		return nil, err
	}
	conn.AddCloseCallback(c.closeConnectionCallback)
	return conn, err
}

func (c *TcpClient) handleRequest(ctx context.Context, conn netpoll.Connection) (err error) {
	reader := conn.Reader()
	len, err := binary.ReadUvarint(reader)
	if err != nil {
		fmt.Printf("%s", err)
		if err == io.EOF {
			conn.Close()
			serevrAddr := ctx.Value("serverAddr")
			if addr, ok := serevrAddr.(string); ok {
				delete(c.connTable, addr)
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
				delete(c.connTable, addr)
			}
		}
	}
	frame, err := Decode(LengthValueBasedCodec, data)
	if err != nil {
		fmt.Printf("%s", err)
		return err
	}
	fmt.Printf("received frame sequence no.: %d", frame.Sequence)
	if _, exists := c.respTable[frame.Sequence]; !exists {
		fmt.Printf("what's wrong? frame sequence not matched with sequence no.: %d", frame.Sequence)
	} else {
		rf := c.respTable[frame.Sequence]
		rf.Add(frame)
	}

	return nil
}

func (c *TcpClient) closeConnectionCallback(conn netpoll.Connection) error {
	fmt.Printf("[%v] connection closed\n", conn.RemoteAddr())
	addr := conn.RemoteAddr()
	conn.Close()
	delete(c.connTable, addr.String())
	return nil
}

func (c *TcpClient) addSeqFuture(seq uint64, rf ResponseFuture) {
	c.mux.Lock()
	defer c.mux.Unlock()

	c.respTable[seq] = rf
}

func (c *TcpClient) cleanupResponseFutures() {
	// 设置定时器，每30秒触发一次扫描
	c.ticker = time.NewTicker(30 * time.Second)
	go func() {
		for {
			<-c.ticker.C
			c.doCleanupRespFutures()
		}
	}()
}

func (c *TcpClient) doCleanupRespFutures() {
	c.mux.Lock()
	defer c.mux.Unlock()
	now := time.Now()
	c.mux.Lock()
	defer c.mux.Unlock()
	for seq, future := range c.respTable {
		if now.Sub(future.Timestamp()) > 30*time.Second {
			// 如果ResponseFuture超过30秒钟
			// 从respTable删除
			delete(c.respTable, seq)
			// 关闭CountDownLatch
			future.Close()
		}
	}
}
