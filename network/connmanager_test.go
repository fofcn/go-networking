package network_test

import (
	"go-networking/network"
	"go-networking/network/util"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestStore_ShouldAddConnection_WhenIdAndConnAreValid(t *testing.T) {
	manager := network.NewConnManager()
	testConn := &network.Conn{
		Connection: nil,
	}
	manager.Store("testID", testConn)
	conn, exists := manager.Load("testID")

	assert.True(t, exists, "Connection should exist")
	assert.Equal(t, testConn, conn, "Should return correct connection")
}

func TestLoad_ShouldReturnConnection_WhenIdExists(t *testing.T) {
	manager := network.NewConnManager()
	testConn := &network.Conn{
		Connection: nil,
	}
	manager.Store("testID", testConn)

	conn, exists := manager.Load("testID")

	assert.True(t, exists, "Connection should exist")
	assert.Equal(t, testConn, conn, "Should return correct connection")
}

func TestPing_ShouldReturnError_WhenIdDoesNotExist(t *testing.T) {
	manager := network.NewConnManager()

	err := manager.Ping("nonexistentID", time.Now().Unix())

	assert.NotNil(t, err, "Should return error")
	assert.Equal(t, "client not found", err.Error(), "Should return correct error")
}

func TestCleanup_ShouldRemoveIdleConnections_WhenTimeoutIsReached(t *testing.T) {
	manager := network.NewConnManager()

	testConn := &network.Conn{
		Connection: nil,
	}
	manager.Store("testID", testConn)

	// 等待超时，然后调用 cleanupNoActiveConn
	time.Sleep(60 * time.Second)

	_, exists := manager.Load("testID")

	assert.False(t, exists, "Connection should have been removed")
}

func TestCleanup_ShouldNotRemoveIdleConnections_WhenPingedInTimeoutIntervalI(t *testing.T) {
	manager := network.NewConnManager()

	testConn := &network.Conn{
		Connection: nil,
	}
	manager.Store("testID", testConn)

	ticker := time.NewTicker(15 * time.Second)
	go func() {
		<-ticker.C
		manager.Ping("1111", time.Now().Unix())
	}()

	// 等待超时，然后调用 cleanupNoActiveConn
	time.Sleep(60 * time.Second)

	_, exists := manager.Load("testID")

	assert.True(t, exists, "Connection should have been existed")
}

func BenchmarkStore(b *testing.B) {
	manager := network.NewConnManager()
	for i := 0; i < b.N; i++ {
		testConn := &network.Conn{Connection: nil}
		manager.Store(util.GetUUIDNoDash(), testConn)
	}
}
