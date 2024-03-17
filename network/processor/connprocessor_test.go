package processor_test

import (
	"fmt"
	"go-networking/network"
	"go-networking/network/codec"
	"go-networking/network/processor"
	"math/big"
	"testing"
	"time"

	"github.com/quintans/toolkit/latch"
	"github.com/stretchr/testify/assert"
)

const (
	serverAddr = "127.0.0.1:8080"
)

func TestProcess_ShouldReturnSuccess_WhenConnIsOK(t *testing.T) {
	network.AddHeaderCodec(network.CONN, &codec.ConnHeaderCodec{})
	network.AddHeaderCodec(network.CONNACK, &codec.ConnAckHeaderCodec{})
	s := StartTcpServer()
	c := StartTcpClient()

	// given
	connHeader := &codec.ConnHeader{
		Timestamp: time.Now().Unix(),
	}
	req := network.NewFrame(network.CONN, connHeader, new(big.Int).SetInt64(1024).Bytes())
	resp, err := c.SendSync(serverAddr, req, 20*time.Second)
	assert.Nil(t, err)
	assert.NotNil(t, resp.Payload)

	s.Stop()
	c.Stop()
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
	tcpServer.AddProcessor(network.CONN, &processor.ConnProcessor{
		TcpServer: tcpServer,
	})
	go func() {
		tcpServer.Start()
		countdownLatch.Done()
	}()

	countdownLatch.Wait()
	defer countdownLatch.Close()

	return tcpServer
}
