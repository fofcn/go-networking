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
)

const (
	serverAddr = "127.0.0.1:8080"
)

func TestSendSyncShouldRecvSuccessResponseWhenConnectedServerAndSendFrameSuccess(t *testing.T) {
	tcpClientConfig := &network.TcpClientConfig{
		Network: "tcp",
		Timeout: time.Duration(time.Duration.Seconds(60)),
	}

	network.AddHeaderCodec(CommandA, &ConnCodecClient{})

	frame := &network.Frame{
		Version: 1,
		CmdType: CommandA,
		Header: &Conn{
			KeyLen: uint32(len("ABC")),
			Key:    "ABC",
		},
		Payload: []byte("Hello world!"),
	}

	countdownLatch := latch.NewCountDownLatch()
	countdownLatch.Add(1)

	startTcpServer(countdownLatch)
	countdownLatch.Wait()
	defer countdownLatch.Close()

	tcpClient := network.NewTcpClient(tcpClientConfig)
	tcpClient.Init()
	tcpClient.Start()

	recvFrame, err := tcpClient.SendSync(serverAddr, frame, 30*time.Second)
	if err != nil {
		log.Printf("error occured when sendSync, %s\n", err)
	} else {
		log.Printf("send sync success, received command type: %d\n", recvFrame.CmdType)
	}

	tcpClient.Stop()
}

func startTcpServer(countdownLatch *latch.CountDownLatch) {
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
		Version:  1,
		CmdType:  CommandA,
		Sequence: packet.Sequence,
		Header: &Conn{
			KeyLen: uint32(len("ABC")),
			Key:    "ABC",
		},
		Payload: []byte("Hello Client"),
	}, nil
}
