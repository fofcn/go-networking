package network

type Packet struct {
}

type Addr struct {
	Host string
	Port string
}

type Conn interface {
}

type RequestInterceptor interface {
	BeforeRequest(remoteAddr string, request *Packet)
	AfterResponse(remoteAddr string, request *Packet, response *Packet)
}

type Processor interface {
	Process(conn *Conn, packet *Packet)
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
	SendSync(addr *Addr, packet *Packet) (*Packet, error)
	SendAsync(addr *Addr, packet *Packet) error
	SendOnce(addr *Addr, packet *Packet) error
	AddProcessor(commandType CommandType, processor Processor)
	AddInterceptor(requestInterceptor RequestInterceptor)
}
