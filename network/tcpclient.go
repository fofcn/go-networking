package network

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"go-networking/log"
	"io"
	"math/big"
	"sync"
	"time"

	"github.com/cloudwego/netpoll"
)

type TcpClientConfig struct {
	Network string
	Timeout time.Duration
}

type HostConn struct {
	id        string
	conn      netpoll.Connection
	seqIncr   *SafeIncrementer32
	key       string
	priKey    big.Int
	timestamp int64
}

type TcpClient struct {
	mux           sync.Mutex
	config        *TcpClientConfig
	hostConnTable map[string]*HostConn
	promiseM      *PromiseM
	ticker        *time.Ticker
	procs         map[CommandType]Processor
	interceptors  []RequestInterceptor
	ctx           context.Context
	cancel        context.CancelFunc
	seqIncr       *SafeIncrementer32
}

func NewTcpClient(config *TcpClientConfig) *TcpClient {
	return &TcpClient{
		config:        config,
		hostConnTable: make(map[string]*HostConn),
		promiseM:      NewPromiseM(),
		ticker:        time.NewTicker(time.Second * 30),
		procs:         make(map[CommandType]Processor, 0),
		interceptors:  make([]RequestInterceptor, 0),
		seqIncr:       NewSafeIncrementer(),
	}
}

func (c *TcpClient) Init() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	c.ctx = ctx
	c.cancel = cancel
	return nil
}

func (c *TcpClient) Start() error {
	c.cleanupResponseFutures()
	return nil
}

func (c *TcpClient) Stop() error {
	c.doCloseConn()
	c.promiseM.CloseRespPromis()
	defer c.cancel()
	c.ticker.Stop()
	return nil
}

func (c *TcpClient) SendSync(serverAddr string, frame *Frame, timeout time.Duration) (*Frame, error) {
	if frame == nil {
		return nil, errors.New("frame is nil")
	}
	frame.Seq = uint64(c.seqIncr.Increment())
	log.Infof("frame auto increment sequence no: %d", frame.Seq)
	rp := NewResponsePromise(frame.Seq, timeout)
	defer rp.Close()
	c.promiseM.AddSeqPromise(frame.Seq, rp)

	err := c.doSendAsync(serverAddr, frame)
	if err != nil {
		return nil, err
	}

	respFrame, err := rp.Wait()
	if err != nil {
		return nil, err
	}

	c.promiseM.DelSeqPromise(frame.Seq)
	return respFrame, nil
}

func (c *TcpClient) SendAsync(serverAddr string, frame *Frame) error {
	frame.Seq = uint64(c.seqIncr.Increment())
	return c.doSendAsync(serverAddr, frame)
}

func (c *TcpClient) doSendAsync(serverAddr string, frame *Frame) error {
	connSeq, err := c.getOrCreateConnection(c.config.Network, serverAddr, c.config.Timeout)
	if err != nil {
		return err
	}
	conn := connSeq.conn
	writer := conn.Writer()
	// encode frame
	bytes, err := Encode(LVBasedCodec, frame)
	if err != nil {
		return err
	}

	cnt, err := writer.WriteBinary(bytes)
	if err != nil || cnt != len(bytes) {
		return errors.New("send failed")
	}

	return writer.Flush()
}

func (c *TcpClient) SendOnce(serverAddr string, packet *Frame) error {
	return errors.New("NotImplemented")
}

func (c *TcpClient) AddProcessor(commandType CommandType, processor Processor) {
	c.mux.Lock()
	defer c.mux.Unlock()
	c.procs[commandType] = processor
}

func (c *TcpClient) AddInterceptor(requestInterceptor RequestInterceptor) {
	c.mux.Lock()
	defer c.mux.Unlock()
	c.interceptors = append(c.interceptors, requestInterceptor)
}

func (c *TcpClient) getOrCreateConnection(network string, serverAddr string, timeout time.Duration) (*HostConn, error) {

	c.mux.Lock()
	defer c.mux.Unlock()

	connSeq, exists := c.hostConnTable[serverAddr]
	if exists {
		return c.recreateConnectionIfNeeded(connSeq, network, serverAddr, timeout)
	}

	return c.doCreateConnection(network, serverAddr, timeout)

}

func (c *TcpClient) recreateConnectionIfNeeded(connSeq *HostConn, network string, serverAddr string, timeout time.Duration) (*HostConn, error) {
	if connSeq.conn.IsActive() {
		return connSeq, nil
	}

	// 关闭不活跃的连接并从表中移除
	err := connSeq.conn.Close()
	if err != nil {
		log.Errorf("Error closing connection: %s", err)
	}
	delete(c.hostConnTable, serverAddr)

	return c.doCreateConnection(network, serverAddr, timeout)
}

func (c *TcpClient) doCreateConnection(network string, serverAddr string, timeout time.Duration) (*HostConn, error) {
	// 尝试创建新连接
	newConn, err := c.createConnection(network, serverAddr, timeout)
	if err != nil {
		return nil, err
	}

	if !newConn.IsActive() {
		return nil, errors.New("connection has not activated after creation")
	}

	newConn.SetOnRequest(c.handleRequest)
	newConn.SetReadTimeout(timeout)
	newConn.SetWriteTimeout(timeout)
	newConn.SetIdleTimeout(timeout * 300)

	newConnSeq := &HostConn{
		conn:    newConn,
		seqIncr: NewSafeIncrementer(),
	}
	c.hostConnTable[serverAddr] = newConnSeq

	return newConnSeq, nil
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
				delete(c.hostConnTable, addr)
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
				delete(c.hostConnTable, addr)
			}
		}
	}
	frame, err := Decode(LVBasedCodec, data)
	if err != nil {
		fmt.Printf("%s", err)
		return err
	}
	log.Infof("client received frame sequence no.: %d", frame.Seq)
	c.promiseM.AddResp(frame)
	return nil
}

func (c *TcpClient) closeConnectionCallback(conn netpoll.Connection) error {
	log.Infof("[%v] connection closed\n", conn.RemoteAddr())
	addr := conn.RemoteAddr()
	conn.Close()
	delete(c.hostConnTable, addr.String())
	return nil
}

func (c *TcpClient) doCloseConn() {
	c.mux.Lock()
	defer c.mux.Unlock()
	for _, connIncr := range c.hostConnTable {
		if connIncr.conn.IsActive() {
			connIncr.conn.Close()
		}
	}
}

func (c *TcpClient) cleanupResponseFutures() {
	// 设置定时器，每30秒触发一次扫描
	c.ticker = time.NewTicker(30 * time.Second)
	go func() {
		for {
			<-c.ticker.C
			c.doCleanupRespPromise()
		}
	}()
}

func (c *TcpClient) doCleanupRespPromise() {
	c.promiseM.CleanupRespPromise()
}
