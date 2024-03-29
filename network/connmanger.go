package network

import (
	"errors"
	"sync"
	"time"
)

func newConnCtx(conn *Conn) *ConnCtx {
	return &ConnCtx{
		Conn:         conn,
		CKey:         "",
		SKey:         "",
		LastPingTime: time.Now().Unix(),
	}
}

func (ctx *ConnCtx) updateCKey(ckey string) {
	ctx.CKey = key
}

func (ctx *ConnCtx) updateSKey(skey string) {
	ctx.SKey = skey
}

func (ctx *ConnCtx) updatePing() {
	ctx.LastPingTime = time.Now().Unix()
}

type ConnManager struct {
	// connKeyTable map[string]*ConnKey
	// string: device-uid
	// connKey: connection, encrypt key, and others
	deviceConnMap *sync.Map
	isStopped     bool
	timeout       time.Duration
	timer         *time.Ticker
}

func NewConnManager() *ConnManager {
	cm := &ConnManager{
		deviceConnMap: &sync.Map{},
		isStopped:     false,
		// todo default timeout is 30s
		timeout: 30 * time.Second,
		timer:   time.NewTicker(30 * time.Second),
	}

	cm.cleanupNoActiveConn()
	return cm
}

func (cm *ConnManager) Store(id string, conn *Conn) {
	cm.deviceConnMap.Store(id, newConnCtx(conn))
}

func (cm *ConnManager) StoreCKey(id string, cKey *string) {
	if value, ok := cm.deviceConnMap.Load(id); ok {
		ctx := value.(*ConnCtx)
		ctx.updateCKey(cKey)
	}
}

func (cm *ConnManager) StoreSKey(id string, sKey *string) {
	if value, ok := cm.deviceConnMap.Load(id); ok {
		ctx := value.(*ConnCtx)
		ctx.updateSKey(sKey)
	}
}

func (cm *ConnManager) Delete(id string) *ConnCtx {
	if value, ok := cm.deviceConnMap.Load(id); ok {
		cm.deviceConnMap.Delete(id)
		return value.(*ConnCtx)
	}

	return nil
}

func (cm *ConnManager) Load(id string) (*Conn, bool) {
	if value, ok := cm.deviceConnMap.Load(id); ok {
		ctx := value.(*ConnCtx)
		return ctx.Conn, true
	}

	return nil, false
}

func (cm *ConnManager) Ping(id string) error {
	if value, ok := cm.deviceConnMap.Load(id); ok {
		if connctx, ok := value.(*ConnCtx); ok {
			connctx.updatePing()
			return nil
		}
	}

	return errors.New("client not found")
}

func (cm *ConnManager) Stop() {
	cm.isStopped = true
	defer cm.timer.Stop()
}

func (cm *ConnManager) cleanupNoActiveConn() {
	go cm.doCleanupNoActiveConn()
}

func (cm *ConnManager) doCleanupNoActiveConn() {
	println("enter cleanup timer")
	for {
		<-cm.timer.C
		if !cm.isStopped {
			println("cleanup timer has triggered")
			cm.deviceConnMap.Range(func(k, v interface{}) bool {
				now := time.Now().Unix()
				if connctx, ok := v.(*ConnCtx); ok {
					if now-connctx.LastPingTime > int64(cm.timeout/time.Second) {
						cm.deviceConnMap.Delete(k)
					}
				}

				return true
			})
		}

	}
}
