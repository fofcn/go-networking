package network_test

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"go-networking/network"
	"log"
	"testing"
	"time"

	"github.com/quintans/toolkit/latch"
	"github.com/stretchr/testify/assert"
)

const (
	serverAddr = "127.0.0.1:8080"
)

func TestSendSyncShouldRecvSuccessResponseWhenConnectedServerAndSendFrameWitiHeaderAndPayload(t *testing.T) {

	network.AddHeaderCodec(CommandA, &ConnCodecClient{})
	startTcpServer()

	frame := &network.Frame{
		Version: 1,
		CmdType: CommandA,
		Header: &Conn{
			KeyLen: uint32(len("ABC")),
			Key:    "ABC",
		},
		Payload: []byte("Hello world!"),
	}

	tcpClient := startTcpClient()
	recvFrame, err := tcpClient.SendSync(serverAddr, frame, 30*time.Second)
	assert.Nil(t, err)
	assert.Equal(t, frame.Seq, recvFrame.Seq)
	fmt.Printf("send sequence no: %d, recv sequence no: %d\n", frame.Seq, recvFrame.Seq)
	tcpClient.Stop()
}

func TestConnectionFailure(t *testing.T) {
	tcpClient := startTcpClient()
	invalidServerAddr := "999.999.999.999:99999" // 这是一个无效的地址和端口
	_, err := tcpClient.SendSync(invalidServerAddr, nil, 30*time.Second)
	assert.NotNil(t, err)
	tcpClient.Stop()
}

func TestHighTrafficStability(t *testing.T) {
	network.AddHeaderCodec(CommandA, &ConnCodecClient{})
	startTcpServer()
	tcpClient := startTcpClient()

	upperLimit := 1000
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

		recvFrame, err := tcpClient.SendSync(serverAddr, frame, 30*time.Second)
		if err != nil {
			t.Errorf("Error occurred on high traffic: %s\n", err)
		} else {
			fmt.Printf("recv frame sequence no: %d", recvFrame.Seq)
		}
	}

	tcpClient.Stop()
}

func startTcpClient() *network.TcpClient {
	tcpClientConfig := &network.TcpClientConfig{
		Network: "tcp",
		Timeout: time.Duration(time.Duration.Seconds(60)),
	}

	tcpClient := network.NewTcpClient(tcpClientConfig)
	tcpClient.Init()
	tcpClient.Start()
	return tcpClient
}

func startTcpServer() {
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
		return
	}
	tcpServer.Init()
	tcpServer.AddProcessor(CommandA, CommandAProcessor{})
	go func() {
		tcpServer.Start()
		countdownLatch.Done()
	}()

	countdownLatch.Wait()
	defer countdownLatch.Close()
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
	log.Println("Command A Process")

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
