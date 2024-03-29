package network

import "github.com/cloudwego/netpoll"

type Addr struct {
	Host string
	Port string
}

type Conn struct {
	Connection netpoll.Connection
}

type ConnCtx struct {
	// real connection
	Conn *Conn
	// client encrypt key
	CKey string
	// server encrypt key
	SKey string
	// last ping time, update by ping command
	LastPingTime int64
}

type RequestInterceptor interface {
	OnRequest(remoteAddr string, request *Frame)
	OnResponse(remoteAddr string, request *Frame, response *Frame)
}

type Processor interface {
	Process(conn *Conn, packet *Frame) (*Frame, error)
}

type Lifecycle interface {
	Init() error
	Start() error
	Stop() error
}

type Server interface {
	Lifecycle
	AddProcessor(cmdType CommandType, process Processor)
	AddInterceptor(requestInterceptor RequestInterceptor)
}

type Client interface {
	Lifecycle
	SendSync(addr *Addr, packet *Frame) (*Frame, error)
	SendAsync(addr *Addr, packet *Frame) error
	SendOnce(addr *Addr, packet *Frame) error
	AddProcessor(commandType CommandType, processor Processor)
	AddInterceptor(requestInterceptor RequestInterceptor)
}
