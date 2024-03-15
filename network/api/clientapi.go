package api

import (
	"go-networking/crypto/dh"
	"go-networking/network"
	"math/big"
	"time"
)

type NasClient interface {
	network.Lifecycle
}

type ConnContext struct {
	id     string
	key    string
	priKey big.Int
}

type TNasClient struct {
	srvAddr     string
	readTimeout time.Duration
	connContext *ConnContext
	tcpClient   *network.TcpClient
}

func NewTNasClient() *TNasClient {
	return &TNasClient{}
}

func (tnc *TNasClient) Init() {
	// load config from yaml file
	// new tcp client from config file
	// tcp client init
	// tcp client start
}

// load config from yaml file
func (tnc *TNasClient) Close() {
	// tcp client stop
}

type ClientApi struct {
	serverAddr  string
	readTimeout time.Duration
	connContext *ConnContext
}

func (c *ClientApi) SendConn(tcpClient network.TcpClient) error {
	var cliPrivateKey, cliPublicKey, err = dh.FastGenDHKP()
	if err != nil {
		return err
	}
	req := network.NewFrame(network.VERSION_1, network.CONN,
		&network.ConnHeader{
			Timestamp: time.Now().Unix(),
		}, cliPublicKey.Bytes(),
	)

	resp, err := tcpClient.SendSync(c.serverAddr, req, c.readTimeout)
	if err != nil {
		return err
	}
	srvPublicKey := new(big.Int).SetBytes(resp.Payload)
	sharedKey := dh.FastGenDHSharedKey(srvPublicKey, cliPrivateKey)
	aesKey := dh.GenAESKeyFromDHKey(sharedKey)
	respHeader := resp.Header.(*network.ConnAckHeader)
	connId := respHeader.Id
	c.connContext = &ConnContext{
		id:     connId,
		key:    string(aesKey),
		priKey: *cliPrivateKey,
	}

	return nil
}
