package codec

import (
	"bytes"
	"encoding/binary"
)

type FileTransfer struct {
	Length   uint32
	FilePath string
}

type FileTransferAck struct {
	FileID    uint32
	FileLen   uint64
	Checksum  uint32
	BlockSize uint32
	ErrorCode uint32
}

type Transfer struct {
	FileID uint32
	Seq    uint32
	Block  []byte
}

// Codec struct for FileTransfer
type FileTransferCodec struct{}

// Codec struct for FileTransferAck
type FileTransferAckCodec struct{}

// Codec struct for Transfer
type TransferCodec struct{}

// Methods for FileTransferCodec
func (ftc *FileTransferCodec) Encode(ft *FileTransfer) ([]byte, error) {
	buf := new(bytes.Buffer)

	if err := binary.Write(buf, binary.BigEndian, ft.Length); err != nil {
		return nil, err
	}

	if err := encodeString(buf, ft.FilePath); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (ftc *FileTransferCodec) Decode(data []byte) (*FileTransfer, error) {
	reader := bytes.NewReader(data)

	ft := &FileTransfer{}

	if err := binary.Read(reader, binary.BigEndian, &ft.Length); err != nil {
		return nil, err
	}

	var err error
	ft.FilePath, err = decodeString(reader)
	if err != nil {
		return nil, err
	}

	return ft, nil
}

// Methods for FileTransferAckCodec
func (ftac *FileTransferAckCodec) Encode(fta *FileTransferAck) ([]byte, error) {
	buf := new(bytes.Buffer)

	if err := binary.Write(buf, binary.BigEndian, fta.FileID); err != nil {
		return nil, err
	}

	if err := binary.Write(buf, binary.BigEndian, fta.FileLen); err != nil {
		return nil, err
	}

	if err := binary.Write(buf, binary.BigEndian, fta.Checksum); err != nil {
		return nil, err
	}

	if err := binary.Write(buf, binary.BigEndian, fta.BlockSize); err != nil {
		return nil, err
	}

	if err := binary.Write(buf, binary.BigEndian, fta.ErrorCode); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (ftac *FileTransferAckCodec) Decode(data []byte) (*FileTransferAck, error) {
	reader := bytes.NewReader(data)

	fta := &FileTransferAck{}

	if err := binary.Read(reader, binary.BigEndian, &fta.FileID); err != nil {
		return nil, err
	}

	if err := binary.Read(reader, binary.BigEndian, &fta.FileLen); err != nil {
		return nil, err
	}

	if err := binary.Read(reader, binary.BigEndian, &fta.Checksum); err != nil {
		return nil, err
	}

	if err := binary.Read(reader, binary.BigEndian, &fta.BlockSize); err != nil {
		return nil, err
	}

	if err := binary.Read(reader, binary.BigEndian, &fta.ErrorCode); err != nil {
		return nil, err
	}

	return fta, nil
}

// Methods for TransferCodec
func (tc *TransferCodec) Encode(tr *Transfer) ([]byte, error) {
	buf := new(bytes.Buffer)

	if err := binary.Write(buf, binary.BigEndian, tr.FileID); err != nil {
		return nil, err
	}

	if err := binary.Write(buf, binary.BigEndian, tr.Seq); err != nil {
		return nil, err
	}

	if err := encodeBytes(buf, tr.Block); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (tc *TransferCodec) Decode(data []byte) (*Transfer, error) {
	reader := bytes.NewReader(data)

	tr := &Transfer{}

	if err := binary.Read(reader, binary.BigEndian, &tr.FileID); err != nil {
		return nil, err
	}

	if err := binary.Read(reader, binary.BigEndian, &tr.Seq); err != nil {
		return nil, err
	}

	var err error
	tr.Block, err = decodeBytes(reader)
	if err != nil {
		return nil, err
	}

	return tr, nil
}

// Helper functions

func encodeString(buf *bytes.Buffer, data string) error {
	if err := binary.Write(buf, binary.BigEndian, uint32(len(data))); err != nil {
		return err
	}
	_, err := buf.WriteString(data)
	return err
}

func decodeString(reader *bytes.Reader) (string, error) {
	var length uint32
	if err := binary.Read(reader, binary.BigEndian, &length); err != nil {
		return "", err
	}
	data := make([]byte, length)
	if _, err := reader.Read(data); err != nil {
		return "", err
	}
	return string(data), nil
}

func encodeBytes(buf *bytes.Buffer, data []byte) error {
	if err := binary.Write(buf, binary.BigEndian, uint32(len(data))); err != nil {
		return err
	}
	_, err := buf.Write(data)
	return err
}

func decodeBytes(reader *bytes.Reader) ([]byte, error) {
	var length uint32
	if err := binary.Read(reader, binary.BigEndian, &length); err != nil {
		return nil, err
	}
	data := make([]byte, length)
	if _, err := reader.Read(data); err != nil {
		return nil, err
	}
	return data, nil
}
