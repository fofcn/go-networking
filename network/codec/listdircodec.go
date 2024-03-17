package codec

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"io"
)

type ListDirHeader struct {
	Id        string
	Timestamp int64
}

type ListDirPayload struct {
	DirPath string
}

type ListDirAckHeader struct {
	StatusCode uint16
}

type ListDirAckPayload struct {
	Files []string
}

func writeLVString(writer io.Writer, data string) error {
	if err := binary.Write(writer, binary.BigEndian, uint16(len(data))); err != nil {
		return err
	}
	_, err := writer.Write([]byte(data))
	return err
}

func readLVString(reader io.Reader) (string, error) {
	var length uint16
	if err := binary.Read(reader, binary.BigEndian, &length); err != nil {
		return "", err
	}
	data := make([]byte, length)
	if _, err := io.ReadFull(reader, data); err != nil {
		return "", err
	}
	return string(data), nil
}

func (h *ListDirHeader) Encode() ([]byte, error) {
	if len(h.Id) != 32 {
		return nil, errors.New("ListDirHeader ID must be 32 bytes long")
	}
	buf := new(bytes.Buffer)
	buf.WriteString(h.Id) // 32 bytes for ID
	binary.Write(buf, binary.BigEndian, h.Timestamp)
	return buf.Bytes(), nil
}

func (h *ListDirHeader) Decode(data []byte) error {
	if len(data) < 40 {
		return errors.New("data is too short to contain a valid ListDirHeader")
	}
	h.Id = string(data[:32])
	buf := bytes.NewReader(data[32:40])
	return binary.Read(buf, binary.BigEndian, &h.Timestamp)
}

func (p *ListDirPayload) Encode() ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := writeLVString(buf, p.DirPath); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (p *ListDirPayload) Decode(data []byte) error {
	dirPath, err := readLVString(bytes.NewReader(data))
	if err != nil {
		return err
	}
	p.DirPath = dirPath
	return nil
}

func (h *ListDirAckHeader) Encode() ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.BigEndian, h.StatusCode); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (h *ListDirAckHeader) Decode(data []byte) error {
	if len(data) != 2 {
		return errors.New("data is not the correct size to contain a valid ListDirAckHeader")
	}
	return binary.Read(bytes.NewReader(data), binary.BigEndian, &h.StatusCode)
}

func (p *ListDirAckPayload) Encode() ([]byte, error) {
	jsonData, err := json.Marshal(p.Files)
	if err != nil {
		return nil, err
	}
	buf := new(bytes.Buffer)
	if err := writeLVString(buf, string(jsonData)); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (p *ListDirAckPayload) Decode(data []byte) error {
	jsonString, err := readLVString(bytes.NewReader(data))
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(jsonString), &p.Files)
}
