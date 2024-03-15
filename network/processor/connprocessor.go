package processor

import "go-networking/network"

type ConnHeader struct {
	Timestamp int64 // Timestamp
}

type ConnPayload struct {
	PublicKey []byte // Client's DH public key
}

type ConnProcessor struct {
}

func (cp *ConnProcessor) Process(conn *network.Conn, packet *network.Frame) (*network.Frame, error) {
	return nil, nil
}
