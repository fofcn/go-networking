package network_test

import (
	"bufio"
	"fmt"
	"go-networking/network"
	"net"
	"testing"
	"time"
)

func TestTcpServer(t *testing.T) {
	ts, _ := network.NewTcpServer()

	err := ts.Init()
	if err != nil {
		t.Errorf("failed to initialize tcp server: %v", err)
		return
	}

	go func() {
		err := ts.Start()
		if err != nil {
			t.Errorf("failed to start tcp server: %v", err)
			return
		}
	}()

	// Wait for the server to start
	time.Sleep(1 * time.Second)

	// Connect to the server
	conn, err := net.Dial("tcp", "localhost:8081")
	if err != nil {
		t.Errorf("failed to connect to tcp server: %v", err)
	}

	// Send a message to the server
	fmt.Fprintf(conn, "Test Message\n")

	// Read the response
	message := bufio.NewReader(conn)
	response, _ := message.ReadString('\n')

	// Check the response
	t.Logf("Received message from server: %s", response)

	// simulate clients
	go func() {
		conn, err := net.Dial("tcp", "localhost:8081")
		if err != nil {
			t.Errorf("failed to connect to server: %v", err)
		}
		defer conn.Close()

		fmt.Fprintf(conn, "hello from client\n")
	}()

	time.Sleep(time.Second)

	ts.Stop()
}
