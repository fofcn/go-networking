package network

import "github.com/cloudwego/netpoll"

type Addr struct {
	Host string
	Port string
}

type Conn struct {
	connection netpoll.Connection
}

type RequestInterceptor interface {
	BeforeRequest(remoteAddr string, request *Frame)
	AfterResponse(remoteAddr string, request *Frame, response *Frame)
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
