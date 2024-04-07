package network_test

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"go-networking/network"
	"runtime"
	"sync"
	"testing"
	"time"

	"github.com/pkg/profile"
	"github.com/quintans/toolkit/latch"
	"github.com/stretchr/testify/assert"
)

const (
	serverAddr = "127.0.0.1:8080"
)

func TestSendSyncShouldRecvSuccessResponseWhenConnectedServerAndSendFrameWitiHeaderAndPayload(t *testing.T) {

	network.AddHeaderCodec(CommandA, &ConnCodecClient{})
	tcpSrv := StartTcpServer()

	frame := &network.Frame{
		Version: 1,
		CmdType: CommandA,
		Header: &Conn{
			KeyLen: uint32(len("ABC")),
			Key:    "ABC",
		},
		Payload: []byte("Hello world!"),
	}

	tcpClient := StartTcpClient()
	recvFrame, err := tcpClient.SendSync(serverAddr, frame, 30*time.Second)
	assert.Nil(t, err)
	assert.Equal(t, frame.Seq, recvFrame.Seq)
	fmt.Printf("send sequence no: %d, recv sequence no: %d\n", frame.Seq, recvFrame.Seq)
	tcpClient.Stop()
	tcpSrv.Stop()
}

func TestConnectionFailure(t *testing.T) {
	tcpClient := StartTcpClient()
	invalidServerAddr := "999.999.999.999:99999" // 这是一个无效的地址和端口
	_, err := tcpClient.SendSync(invalidServerAddr, &network.Frame{}, 30*time.Second)
	assert.NotNil(t, err)
	tcpClient.Stop()
}

func TestHighTrafficStability(t *testing.T) {
	runtime.GOMAXPROCS(32)
	network.AddHeaderCodec(CommandA, &ConnCodecClient{})
	tcpServer := StartTcpServer()
	tcpClient := StartTcpClient()

	wg := sync.WaitGroup{}
	upperLimit := 25
	for i := 0; i < upperLimit; i++ {
		frame := &network.Frame{
			Version: 1,
			CmdType: CommandA,
			Header: &Conn{
				KeyLen: uint32(len("ABC")),
				Key:    "ABC",
			},
			Payload: []byte("Hello world!"),
		}

		wg.Add(1)
		go func() {
			recvFrame, err := tcpClient.SendSync(serverAddr, frame, 30*time.Second)
			if err != nil {
				t.Errorf("Error occurred on high traffic: %s\n", err)
			} else {
				fmt.Printf("recv frame sequence no: %d\n", recvFrame.Seq)
			}

			wg.Done()
		}()

	}

	wg.Wait()

	tcpClient.Stop()
	tcpServer.Stop()
}

func BenchmarkHighTrafficStability(b *testing.B) {
	// 启动性能分析（可选）
	defer profile.Start().Stop()

	network.AddHeaderCodec(CommandA, &ConnCodecClient{})
	tcpServer := StartTcpServer()
	tcpClient := StartTcpClient()

	b.ResetTimer() // 重置基准计时器

	for i := 0; i < b.N; i++ { // 使用 b.N 作为循环次数
		frame := &network.Frame{
			Version: 1,
			CmdType: CommandA,
			Header: &Conn{
				KeyLen: uint32(len("ABC")),
				Key:    "ABC",
			},
			Payload: []byte("Hello world!"),
		}

		_, _ = tcpClient.SendSync(serverAddr, frame, 30*time.Second) // 忽略返回值，仅关注发送操作
	}

	// 停止 TCP 客户端和服务器（通常在基准测试中不会执行，因为每次迭代都会重新启动它们）
	tcpClient.Stop()
	tcpServer.Stop()
}

func BenchmarkHighTrafficConcurrency(b *testing.B) {
	network.AddHeaderCodec(CommandA, &ConnCodecClient{})
	tcpServer := StartTcpServer()
	tcpClient := StartTcpClient()

	b.ResetTimer() // 重置基准计时器

	for i := 0; i < b.N; i++ {
		frame := &network.Frame{
			Version: 1,
			CmdType: CommandA,
			Header: &Conn{
				KeyLen: uint32(len("ABC")),
				Key:    "ABC",
			},
			Payload: []byte("Hello world!"),
		}

		// 创建并发 Goroutine 发送请求
		var wg sync.WaitGroup
		wg.Add(1)

		go func() {
			defer wg.Done()

			_, _ = tcpClient.SendSync(serverAddr, frame, 30*time.Second) // 忽略返回值，仅关注发送操作
		}()

		wg.Wait() // 等待当前并发请求完成
	}

	// 停止 TCP 客户端和服务器（通常在基准测试中不会执行，因为每次迭代都会重新启动它们）
	tcpClient.Stop()
	tcpServer.Stop()
}

func StartTcpClient() *network.TcpClient {
	tcpClientConfig := &network.TcpClientConfig{
		Network: "tcp",
		Timeout: time.Duration(time.Duration.Seconds(60)),
	}

	tcpClient := network.NewTcpClient(tcpClientConfig)
	tcpClient.Init()
	tcpClient.Start()
	return tcpClient
}

func StartTcpServer() *network.TcpServer {
	countdownLatch := latch.NewCountDownLatch()
	countdownLatch.Add(1)
	addr := network.Addr{
		Host: "127.0.0.1",
		Port: "8080",
	}
	tcpServerConfig := &network.TcpServerConfig{
		Network: "tcp",
		Addr:    addr,
	}

	tcpServer, err := network.NewTcpServer(tcpServerConfig)
	if err != nil {
		fmt.Println("failed to create tcp server")
		return nil
	}
	tcpServer.Init()
	tcpServer.AddProcessor(CommandA, CommandAProcessor{})
	go func() {
		tcpServer.Start()
		countdownLatch.Done()
	}()

	countdownLatch.Wait()
	defer countdownLatch.Close()

	return tcpServer
}

type ConnClient struct {
	KeyLen uint32
	Key    string
}

type ConnCodecClient struct {
}

func (codec ConnCodecClient) Encode(header interface{}) ([]byte, error) {
	if conn, ok := header.(*Conn); ok {
		buf := new(bytes.Buffer)
		buf.Write(network.EncodeInteger(uint64(conn.KeyLen)))
		buf.Write([]byte(conn.Key))
		return buf.Bytes(), nil
	} else {
		// todo
		return nil, errors.New("error occured when try to encode, invalid type of Conn")
	}
}
func (codec ConnCodecClient) Decode(data []byte) (interface{}, error) {
	buf := bytes.NewReader(data)
	keyLen, err := binary.ReadUvarint(buf)
	if err != nil {
		return nil, err
	}

	conn := &Conn{}
	conn.KeyLen = uint32(keyLen)

	keyData := make([]byte, keyLen)

	if _, err := buf.Read(keyData); err != nil {
		return nil, errors.New("failed to read payload")
	}

	conn.Key = string(keyData)
	return conn, nil
}

type CommandAProcessor struct{}

func (cmdProcessor CommandAProcessor) Process(conn *network.Conn, packet *network.Frame) (*network.Frame, error) {
	return &network.Frame{
		Version: 1,
		CmdType: CommandA,
		Seq:     packet.Seq,
		Header: &Conn{
			KeyLen: uint32(len("ABC")),
			Key:    "ABC",
		},
		Payload: []byte("Hello Client"),
	}, nil
}
