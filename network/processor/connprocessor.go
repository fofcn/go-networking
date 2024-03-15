package processor

import (
	"go-networking/crypto/dh"
	"go-networking/network"
	"math/big"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
)

// 用于保护connTable的读写锁
var connTableLock sync.RWMutex
var srvPrivateKey, srvPublicKey, _ = dh.FastGenDHKP()

type ConnProcessor struct {
	TcpServer *network.TcpServer
}

// 实现连接建立
// 获取DH公钥
// 获取自己的私钥后计算出AES加密KEY
// 生成一个连接ID
// 将连接添加到connTable
// 回复连接ID,连接ID写入到ConnHeader的Id字段中
// 获取自己的公钥回复给客户端，公钥写入到ConnPayload的PublicKey字段中
// 使用FastGenDHKP生成DH密钥对
// 使用FastGenDHSharedKey计算共享密钥

func (cp *ConnProcessor) Process(conn *network.Conn, frame *network.Frame) (*network.Frame, error) {
	// 解析客户端DH公钥
	clientDHKey := new(big.Int).SetBytes(frame.Payload)

	// 计算共享秘钥
	sharedKey := dh.FastGenDHSharedKey(clientDHKey, srvPrivateKey)

	// 从共享秘钥生成AES密钥
	aesKey := dh.GenAESKeyFromDHKey(sharedKey)

	// 假设我们有一个方法来生成连接ID并确保其唯一性
	connID := GetUUID()

	// 添加连接到连接表
	connTableLock.Lock()
	// 记录aesKey到连接上下文中，用于之后数据的加解密
	cp.TcpServer.AddConnKey(connID, &network.ConnKey{
		Conn: conn,
		Key:  string(aesKey),
	})
	connTableLock.Unlock()

	// 准备回复客户端的数据包
	respHeader := &network.ConnAckHeader{
		Id:        connID, // 此处需要将connID从uint64转换为字符串，如使用strconv.Itoa
		Timestamp: time.Now().Unix(),
	}

	responseFrame := &network.Frame{
		Version: uint16(network.VERSION_1),
		CmdType: network.CONNACK,
		Seq:     frame.Seq,
		Header:  respHeader,
		Payload: srvPublicKey.Bytes(),
	}

	return responseFrame, nil
}

func GetUUID() string {
	// 生成UUID并去掉所有短横线
	return strings.ReplaceAll(uuid.New().String(), "-", "")
}