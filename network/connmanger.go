package network

import (
	"sync"
	"time"
)

const (
	// 定义超时时长为1小时
	HourInSeconds = 60 * 60
)

type ConnKey struct {
	Conn      *Conn
	Key       string
	Timestamp int64
}

type ConnManager struct {
	mux          sync.Mutex
	connKeyTable map[string]*ConnKey
	isStopped    bool
}

func NewConnManager() *ConnManager {
	cm := &ConnManager{
		connKeyTable: make(map[string]*ConnKey),
		isStopped:    false,
	}

	cm.cleanupNoActiveConn()
	return cm
}

func (cm *ConnManager) AddConn(id string, connKey *ConnKey) {
	cm.mux.Lock()
	defer cm.mux.Unlock()
	cm.connKeyTable[id] = connKey
}

func (cm *ConnManager) Ping(id string, ts int64) {
	if connKey, exists := cm.connKeyTable[id]; exists {
		connKey.Timestamp = ts
	}
}

func (cm *ConnManager) Stop() {
	cm.isStopped = true
}

func (cm *ConnManager) cleanupNoActiveConn() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for cm.isStopped {
		// 这里就是典型判断Channel是否有事件发生，如果发生则执行cm.cleanupNoAciveConn
		<-ticker.C
		cm.doCleanupNoActiveConn()
	}
}

func (cm *ConnManager) doCleanupNoActiveConn() {
	// 锁定connKeyTable的操作
	cm.mux.Lock()
	defer cm.mux.Unlock()

	// 获取当前时间的时间戳
	now := time.Now().Unix()

	// 遍历connKeyTable
	for id, connKey := range cm.connKeyTable {
		// 如果连接的Timestamp小于(当前时间减去1小时)，则删除该连接
		if now-connKey.Timestamp > HourInSeconds {
			delete(cm.connKeyTable, id)
		}
	}
}
