package network

import (
	"errors"
	"sync"
	"time"
)

// newConnCtx 创建一个新的ConnCtx实例。
// conn: 表示连接的Conn对象。
// 返回值: 初始化后的ConnCtx指针。
func newConnCtx(conn *Conn) *ConnCtx {
	return &ConnCtx{
		Conn:         conn,
		CKey:         "",
		SKey:         "",
		LastPingTime: time.Now().Unix(),
	}
}

// updateCKey 更新ConnCtx实例的CKey字段。
// ckey: 要更新的CKey字符串。
func (ctx *ConnCtx) updateCKey(ckey string) {
	ctx.CKey = ckey
}

// updateSKey 更新ConnCtx实例的SKey字段。
// skey: 要更新的SKey字符串。
func (ctx *ConnCtx) updateSKey(skey string) {
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
	isStopped     bool
	timeout       time.Duration // 连接超时时间
	timer         *time.Ticker  // 定时器，用于定期清理无活跃连接
}

// NewConnManager 创建并初始化一个新的ConnManager实例。
// 返回值: 初始化后的ConnManager指针。
func NewConnManager() *ConnManager {
	cm := &ConnManager{
		deviceConnMap: &sync.Map{},
		isStopped:     false,
		timeout:       30 * time.Second,
		timer:         time.NewTicker(30 * time.Second),
	}

	cm.cleanupNoActiveConn()
	return cm
}

// Store 将连接存储到设备连接映射中。
// id: 设备UID。
// conn: 要存储的连接对象。
func (cm *ConnManager) Store(id string, conn *Conn) {
	cm.deviceConnMap.Store(id, newConnCtx(conn))
}

// StoreCKey 更新指定设备的CKey。
// id: 设备UID。
// cKey: 要更新的CKey。
func (cm *ConnManager) StoreCKey(id string, cKey string) {
	if value, ok := cm.deviceConnMap.Load(id); ok {
		ctx := value.(*ConnCtx)
		ctx.updateCKey(cKey)
	}
}

// StoreSKey 更新指定设备的SKey。
// id: 设备UID。
// sKey: 要更新的SKey。
func (cm *ConnManager) StoreSKey(id string, sKey string) {
	if value, ok := cm.deviceConnMap.Load(id); ok {
		ctx := value.(*ConnCtx)
		ctx.updateSKey(sKey)
	}
}

// Delete 从设备连接映射中删除指定的设备连接。
// id: 设备UID。
// 返回值: 被删除的ConnCtx实例，如果未找到则返回nil。
func (cm *ConnManager) Delete(id string) *ConnCtx {
	if value, ok := cm.deviceConnMap.Load(id); ok {
		cm.deviceConnMap.Delete(id)
		return value.(*ConnCtx)
	}

	return nil
}

// Load 根据设备UID加载对应的连接。
// id: 设备UID。
// 返回值: 对应的连接对象和一个布尔值，表示是否成功找到连接。
func (cm *ConnManager) Load(id string) (*Conn, bool) {
	if value, ok := cm.deviceConnMap.Load(id); ok {
		ctx := value.(*ConnCtx)
		return ctx.Conn, true
	}

	return nil, false
}

// Ping 根据设备UID标记该设备连接为活跃。
// id: 设备UID。
// 返回值: 错误对象，如果设备未找到则返回错误。
func (cm *ConnManager) Ping(id string) error {
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
	cm.isStopped = true
	defer cm.timer.Stop()
}

// cleanupNoActiveConn 启动一个goroutine定期清理长时间未活跃的连接。
func (cm *ConnManager) cleanupNoActiveConn() {
	go cm.doCleanupNoActiveConn()
}

// doCleanupNoActiveConn 实际执行清理无活跃连接的任务。
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
