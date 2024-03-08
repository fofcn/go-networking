package network_test

import (
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"go-networking/network"
	"testing"
	"time"

	"golang.org/x/sync/semaphore"
)

const (
	serverAddr = "127.0.0.1:8080"
)

func TestSendSyncShouldRecvSuccessResponseWhenConnectedServerAndSendFrameSuccess(t *testing.T) {
	tcpClientConfig := &network.TcpClientConfig{
		Network: "tcp",
		Timeout: time.Duration(time.Duration.Seconds(60)),
	}

	network.AddCodec(CommandA, &ConnCodecClient{})

	frame := &network.Frame{
		Version: 1,
		CmdType: CommandA,
		Payload: []byte("Hello world!"),
	}

	sem := semaphore.NewWeighted(1)

	go startTcpServer(sem)
	sem.Acquire(context.Background(), 1)

	tcpClient := network.NewTcpClient(tcpClientConfig)
	tcpClient.Init()
	tcpClient.Start()

	tcpClient.SendSync(serverAddr, frame)
}

func startTcpServer(sem *semaphore.Weighted) {
	defer sem.Release(1)
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
	tcpServer.Start()

}

type ConnClient struct {
	KeyLen uint32
	Key    string
}

type ConnCodecClient struct {
}

func (codec ConnCodecClient) Encode(header interface{}) ([]byte, error) {
	if conn, ok := header.(Conn); ok {
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
