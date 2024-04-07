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
	seqIncr   SafeIncrementer32
	key       string
	priKey    big.Int
	timestamp int64
}

type TcpClient struct {
	mux           sync.Mutex
	config        *TcpClientConfig
	hostConnTable map[string]*HostConn
	rpTable       map[uint64]ResponsePromise
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
		rpTable:       make(map[uint64]ResponsePromise),
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
	c.closeRespPromis()
	defer c.cancel()
	c.ticker.Stop()
	return nil
}

func (c *TcpClient) SendSync(serverAddr string, frame *Frame, timeout time.Duration) (*Frame, error) {
	frame.Seq = uint64(c.seqIncr.Increment())
	log.Infof("frame auto increment sequence no: %d", frame.Seq)
	rp := NewResponsePromise(frame.Seq, timeout)
	defer rp.Close()
	c.addSeqPromise(frame.Seq, rp)

	err := c.doSendAsync(serverAddr, frame)
	if err != nil {
		return nil, err
	}

	respFrame, err := rp.Wait()
	if err != nil {
		return nil, err
	}

	c.delSeqPromise(frame.Seq)
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
		if connSeq.conn.IsActive() {
			return connSeq, nil
		} else {
			err := connSeq.conn.Close()
			log.Infof("connection was not active, close also occured error, please check the error: %s", err)
			delete(c.hostConnTable, serverAddr)
		}
	}

	conn, err := c.createConnection(network, serverAddr, timeout)
	if err != nil {
		return nil, err
	}

	if !conn.IsActive() {
		return nil, errors.New("connection have not actived after created")
	}

	conn.SetOnRequest(c.handleRequest)

	connSeq = &HostConn{
		conn:    conn,
		seqIncr: *NewSafeIncrementer(),
	}
	c.hostConnTable[serverAddr] = connSeq

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
	c.addResp(frame)
	return nil
}

func (c *TcpClient) addResp(frame *Frame) {
	c.mux.Lock()
	defer c.mux.Unlock()
	if _, exists := c.rpTable[frame.Seq]; exists {
		rp := c.rpTable[frame.Seq]
		rp.Add(frame)
	} else {
		log.Infof("what's wrong? frame sequence not matched with sequence no.: %d", frame.Seq)
	}
}

func (c *TcpClient) closeConnectionCallback(conn netpoll.Connection) error {
	log.Infof("[%v] connection closed\n", conn.RemoteAddr())
	addr := conn.RemoteAddr()
	conn.Close()
	delete(c.hostConnTable, addr.String())
	return nil
}

func (c *TcpClient) addSeqPromise(seq uint64, rp ResponsePromise) {
	c.mux.Lock()
	defer c.mux.Unlock()

	c.rpTable[seq] = rp
}

func (c *TcpClient) delSeqPromise(seq uint64) {
	c.mux.Lock()
	defer c.mux.Unlock()
	if future, exists := c.rpTable[seq]; exists {
		future.Close()
		delete(c.rpTable, seq)
	}
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

func (c *TcpClient) closeRespPromis() {
	c.mux.Lock()
	defer c.mux.Unlock()

	for seq, future := range c.rpTable {
		delete(c.rpTable, seq)
		// 关闭CountDownLatch
		future.Close()
	}
}

func (c *TcpClient) doCleanupRespPromise() {
	now := time.Now()
	c.mux.Lock()
	defer c.mux.Unlock()
	for seq, future := range c.rpTable {
		if now.Sub(future.Timestamp()) > 30*time.Second {
			// 如果ResponseFuture超过30秒钟
			// 从respTable删除
			delete(c.rpTable, seq)
			// 关闭CountDownLatch
			future.Close()
		}
	}
}
