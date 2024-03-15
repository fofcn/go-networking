# go-networking
This is a go networking examples including TCP and HTTP. HTTP will use gin framework, TCP will use net-poll framework.


# Protocol
Frame
```go
type Frame struct = {
    Version uint16
    CmdType uint32
    Seq     uint64
    HLen    uint16
    Header  interface{}
    Payload []byte
}

```

## Introduction
1. Version: Protocol Version
2. CmdType: Command Type, used for match frame handler
3. Seq: Frame Sequence, used for request response matching
4. HLen: Varint Header Length, indicates Header length
5. Header: Header data
6. Payload: Actual Data, optional. A control frame may not contain a payload

## Cmd Type
1. CONN
2. CONNACK

CONN and CONNACK are used for end-to-end encryption. The client will send its DH public key to the server, and likewise, the server will also send its DH public key to the client.

CONN struct and Frame:
```go
type ConnHeader struct {
    Version    uint16 // Protocol version
    Timestamp  int64  // Timestamp
}

type ConnPayload struct {
    PublicKey []byte // Client's DH public key
}

// When constructing the actual frame:
connFrame := Frame{
    Version: 1,
    CmdType: CONN, // Assuming CONN is defined as a constant somewhere
    Seq:     ...,
    HLen:    ...,
    Header:  ConnHeader{...},
    Payload: ConnPayload{...},
}
```

CONNACK struct and frame:
```go
type ConnAckHeader struct {
    StatusCode uint16 // Status of the connection
}

type ConnAckPayload struct {
    PublicKey []byte // Server's DH public key
}

// When constructing the actual frame:
connAckFrame := Frame{
    Version: 1,
    CmdType: CONNACK, // Assuming CONNACK is defined as a constant somewhere
    Seq:     ...,
    HLen:    ...,
    Header:  ConnAckHeader{...},
    Payload: ConnAckPayload{...},
}
```

3. PING
4. PONG
PING and PONG are used for maintaining the session connection. The client sends a PING to the server, and the server responds with a PONG to the client. If the client does not send a PING to the server within a certain period, the server will close the connection. Similarly, if the client does not receive a PONG from the server within a certain period, the client will close the connection.

PING struct and frame:
```go
type PingHeader struct {
    Timestamp int64 // Timestamp
}

// No payload needed for PING

// Constructing the PING frame
pingFrame := Frame{
    Version: 1,
    CmdType: PING, // Assuming PING is defined as a constant somewhere
    Seq:     ...,
    HLen:    ...,
    Header:  PingHeader{...},
    Payload: nil, // No payload for PING
}
```
PONG struct and frame:
```go
type PongHeader struct {
    Timestamp int64 // Timestamp
}

// No payload needed for PONG

// Constructing the PONG frame
pongFrame := Frame{
    Version: 1,
    CmdType: PONG, // Assuming PONG is defined as a constant somewhere
    Seq:     ...,
    HLen:    ...,
    Header:  PongHeader{...},
    Payload: nil, // No payload for PONG
}
```

5. CLOSE
6. CLOSEACK
CLOSE and CLOSEACK are used to close the connection. Both the client and the server can send a CLOSE frame. If the client sends a CLOSE frame to the server, the server will check if there are any frames yet to be sent to the client. If there are, the server will continue sending these frames to the client. Once completed, the server responds with a CLOSEACK to the client. Upon receiving the CLOSEACK, the client will close the connection.

CLOSE struct and frame:
```go
type CloseHeader struct {
    Reason string // Reason for closing the connection
}

// No payload needed for CLOSE

// Constructing the CLOSE frame
closeFrame := Frame{
    Version: 1,
    CmdType: CLOSE, // Assuming CLOSE is defined as a constant somewhere
    Seq:     ...,
    HLen:    ...,
    Header:  CloseHeader{...},
    Payload: nil, // No payload for CLOSE
}
```

CLOSEACK struct and frame:
```go
type CloseAckHeader struct {
    StatusCode uint16 // Status code of the close operation, e.g., success, failure, etc.
    Details string   // Additional details or reason of the status
}

// No payload needed for CLOSEACK

// Constructing the CLOSEACK frame
closeAckFrame := Frame{
    Version: 1,
    CmdType: CLOSEACK, // Assuming CLOSEACK is defined as a constant somewhere
    Seq:     ...,
    HLen:    ...,
    Header:  CloseAckHeader{...},
    Payload: nil, // No payload for CLOSEACK
}
```

7. LISTDIR
8. LISTDIRACK
LISTDIR and LISTDIRACK are used for listing directory and files.

LISTDIR struct and frame:
```go
type ListDirHeader struct {
    Timestamp int64 // Timestamp
}

type ListDirPayload struct {
    DirPath string // The directory path to be listed
}

// Constructing the LISTDIR frame
listDirFrame := Frame{
    Version: 1,
    CmdType: LISTDIR, // Assuming LISTDIR is defined as a constant somewhere
    Seq:     ...,
    HLen:    ...,
    Header:  ListDirHeader{...},
    Payload: ListDirPayload{...}, 
}
```

LISTDIRACK struct and frame:
```go
type ListDirAckHeader struct {
    StatusCode uint16 // Status code of the operation, e.g., success, failure, etc.
}

type ListDirAckPayload struct {
    Files []string // The files in the directory
}

// Constructing the LISTDIRACK frame
listDirAckFrame := Frame{
    Version: 1,
    CmdType: LISTDIRACK, // Assuming LISTDIRACK is defined as a constant somewhere
    Seq:     ...,
    HLen:    ...,
    Header:  ListDirAckHeader{...},
    Payload: ListDirAckPayload{...},
}
```

9. FILETRANSFER
10. FILETRANSFERACK
11. TRANSFER
Client creates a new connection for file transfer purpose. 

FILETRANSFER and FILETRANSFERACK are used for negotiate file transfer.
Client send FILETRANSFER frame to server, and server send FILETRANSFERACK to client.
FILETRANSFER frame contains filename and path, server gets the metadata of the file, and send ACK to client, ACK contains file name, file length, and block size which represents size of each part. server will encode the file name with an integer number.

TRANSFER is used for server side to transfer file data to client
```go
FILETRANSFER = {
    length uint32,
    filepath string,
}
```

```go
FILETRANSFERACK = {
    fileId uint32,
    fileLen uint64,
    checksum utin32,
    blockSize uint32,
    errorCode uint32,
}
```

```go
TRANSFER = {
    fileId uint32,
    seq uint32,
    block []byte,
}

```
