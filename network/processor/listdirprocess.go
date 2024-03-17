package processor

import (
	"encoding/json"
	"go-networking/log"
	"go-networking/network"
	"go-networking/network/codec"
	"os"
)

type ListdireProcessor struct {
	TcpClient *network.TcpClient
}

func (lp *ListdireProcessor) Process(conn *network.Conn, frame *network.Frame) (*network.Frame, error) {
	header := frame.Header.(*codec.ListDirHeader)

	// todo ping should add to interceptor
	// lp.TcpServer.CManager.Ping(header.Id, header.Timestamp)

	listdir := string(frame.Payload)

	entries, err := os.ReadDir(listdir)
	if err != nil {
		log.Errorf("failed to read directory: %s", err)
		return &network.NewFrame(
			network.LISTDIRACK,
			&codec.ListDirAckHeader{
				StatusCode: 1, // todo need to define status code, 0 indicates success, otherwise failure.
			},
			[]byte("failed to read directory or directory is not existing"), // todo define error message
		)
	}

	var filenames []string
	for _, entry := range entries {
		filenames = append(filenames, e.Name())
	}

	filenamestr := json.Marshal(filenames)

	return &network.NewFrame(
		network.LISTDIRACK,
		&codec.ListDirAckHeader{
			StatusCode: 0, // todo need to define status code, 0 indicates success, otherwise failure.
		},
		[]byte(filenamestr), // todo define error message
	)
}
