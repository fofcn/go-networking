package network

import (
	"errors"
	"go-networking/network/codec"
	"sync"
	"sync/atomic"
	"time"
)

// newConnCtx 创建一个新的ConnCtx实例。
// conn: 表示连接的Conn对象。
// 返回值: 初始化后的ConnCtx指针。
func newConnCtx(conn *Conn, key []byte) *ConnCtx {
	return &ConnCtx{
		Conn:         conn,
		CKey:         key,
		SKey:         key,
		LastPingTime: time.Now().Unix(),
	}
}

// updateCKey 更新ConnCtx实例的CKey字段。
// ckey: 要更新的CKey字符串。
func (ctx *ConnCtx) updateCKey(ckey []byte) {
	ctx.CKey = ckey
}

// updateSKey 更新ConnCtx实例的SKey字段。
// skey: 要更新的SKey字符串。
func (ctx *ConnCtx) updateSKey(skey []byte) {
	ctx.SKey = skey
}

// updatePing 更新ConnCtx实例的LastPingTime字段为当前时间。
func (ctx *ConnCtx) updatePing() {
	ctx.LastPingTime = time.Now().Unix()
}

// ConnManager 是用于管理连接的结构体。
type ConnManager struct {
	// deviceConnMap 用于存储设备连接信息的映射。
	// key: 设备UID。
	// value: ConnCtx实例，包含连接及相关密钥信息。
	deviceConnMap *sync.Map
	isStopped     atomic.Uint32
	timeout       time.Duration // 连接超时时间
	timer         *time.Ticker  // 定时器，用于定期清理无活跃连接
}

// NewConnManager 创建并初始化一个新的ConnManager实例。
func NewConnManager() *ConnManager {
	cm := &ConnManager{
		deviceConnMap: &sync.Map{},
		timeout:       30 * time.Second,
		timer:         time.NewTicker(30 * time.Second),
	}

	cm.isStopped.Store(0)
	go cm.cleanupNoActiveConn()
	return cm
}

// Store 将连接存储到设备连接映射中。
func (cm *ConnManager) Store(id string, conn *Conn, key []byte) {
	cm.deviceConnMap.Store(id, newConnCtx(conn, key))
}

// StoreCKey 更新指定设备的CKey。
func (cm *ConnManager) StoreCKey(id string, cKey []byte) {
	if value, ok := cm.deviceConnMap.Load(id); ok {
		ctx := value.(*ConnCtx)
		ctx.updateCKey(cKey)
	}
}

// StoreSKey 更新指定设备的SKey。
func (cm *ConnManager) StoreSKey(id string, sKey []byte) {
	if value, ok := cm.deviceConnMap.Load(id); ok {
		ctx := value.(*ConnCtx)
		ctx.updateSKey(sKey)
	}
}

// Delete 从设备连接映射中删除指定的设备连接。
// 返回值: 被删除的ConnCtx实例，如果未找到则返回nil。
func (cm *ConnManager) Delete(id string) *ConnCtx {
	if value, ok := cm.deviceConnMap.Load(id); ok {
		cm.deviceConnMap.Delete(id)
		return value.(*ConnCtx)
	}

	return nil
}

// Load 根据设备UID加载对应的连接。
func (cm *ConnManager) Load(id string) (*Conn, bool) {
	if value, ok := cm.deviceConnMap.Load(id); ok {
		ctx := value.(*ConnCtx)
		return ctx.Conn, true
	}

	return nil, false
}

// Ping 根据设备UID标记该设备连接为活跃。
func (cm *ConnManager) Ping(id string, ts int64) error {
	if time.Now().Unix()-ts > int64(cm.timeout/time.Second) {
		return codec.Invalid_Ping_Frame
	}

	if value, ok := cm.deviceConnMap.Load(id); ok {
		if connctx, ok := value.(*ConnCtx); ok {
			connctx.updatePing()
			return nil
		}
	}

	return errors.New("client not found")
}

// Stop 停止ConnManager的定时清理任务。
func (cm *ConnManager) Stop() {
	cm.isStopped.Store(1)
	cm.timer.Stop()
}

// cleanupNoActiveConn 启动一个goroutine定期清理长时间未活跃的连接。
func (cm *ConnManager) cleanupNoActiveConn() {
	for {
		select {
		case <-cm.timer.C:
			if cm.isStopped.Load() == 0 {
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
}
