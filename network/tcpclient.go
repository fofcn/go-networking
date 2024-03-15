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

type ConnSeq struct {
	conn    netpoll.Connection
	seqIncr SafeIncrementer32
}

type TcpClient struct {
	mux       sync.Mutex
	config    *TcpClientConfig
	connTable map[string]*ConnSeq
	respTable map[uint64]ResponsePromise
	ticker    *time.Ticker
}

func NewTcpClient(config *TcpClientConfig) *TcpClient {
	return &TcpClient{
		config:    config,
		connTable: make(map[string]*ConnSeq),
		respTable: make(map[uint64]ResponsePromise),
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
	for _, connIncr := range c.connTable {
		if connIncr.conn.IsActive() {
			connIncr.conn.Close()
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
	connSeq, err := c.getOrCreateConnection(c.config.Network, serverAddr, c.config.Timeout)
	if err != nil {
		return nil, err
	}
	conn := connSeq.conn
	// 创建一个超时上下文
	// 设置30秒超时
	// 在完成后释放资源
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	var key contextKey = "serverAddr"
	context.WithValue(ctx, key, serverAddr)
	conn.SetOnRequest(c.handleRequest)

	writer := conn.Writer()
	// encode frame
	frame.Seq = uint64(connSeq.seqIncr.Increment())
	bytes, err := Encode(LVBasedCodec, frame)
	if err != nil {
		return nil, err
	}

	rp := NewResponsePromise(frame.Seq, timeout)
	defer rp.Close()
	c.addSeqFuture(frame.Seq, rp)

	cnt, err := writer.WriteBinary(bytes)
	if err != nil || cnt != len(bytes) {
		return nil, errors.New("send failed")
	}

	err = writer.Flush()
	if err != nil {
		return nil, err
	}

	respFrame, err := rp.Wait()
	if err != nil {
		return nil, err
	}
	delete(c.respTable, frame.Seq)
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

func (c *TcpClient) getOrCreateConnection(network string, serverAddr string, timeout time.Duration) (*ConnSeq, error) {
	c.mux.Lock()
	defer c.mux.Unlock()
	connSeq, exists := c.connTable[serverAddr]
	if exists {
		if connSeq.conn.IsActive() {
			return connSeq, nil
		} else {
			err := connSeq.conn.Close()
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

	connSeq = &ConnSeq{
		conn:    conn,
		seqIncr: *NewSafeIncrementer(),
	}
	c.connTable[serverAddr] = connSeq

	return connSeq, nil
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
	frame, err := Decode(LVBasedCodec, data)
	if err != nil {
		fmt.Printf("%s", err)
		return err
	}
	fmt.Printf("received frame sequence no.: %d", frame.Seq)
	if _, exists := c.respTable[frame.Seq]; !exists {
		fmt.Printf("what's wrong? frame sequence not matched with sequence no.: %d", frame.Seq)
	} else {
		rp := c.respTable[frame.Seq]
		rp.Add(frame)
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

func (c *TcpClient) addSeqFuture(seq uint64, rp ResponsePromise) {
	c.mux.Lock()
	defer c.mux.Unlock()

	c.respTable[seq] = rp
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
