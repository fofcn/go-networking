package processor

type ConnAckHeader struct {
	StatusCode uint16 // Status of the connection
}

type ConnAckPayload struct {
	PublicKey []byte // Server's DH public key
}
