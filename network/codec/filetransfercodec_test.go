package codec_test

import (
	"bytes"
	"encoding/binary"
	"go-networking/network/codec"
	"reflect"
	"testing"
)

func TestFileTransferCodec_Encode_ShouldReturnBytes_WhenGivenFileTransfer(t *testing.T) {
	ftc := codec.FileTransferCodec{}
	ft := codec.FileTransfer{
		Length:   4,
		FilePath: "/path/to/file",
	}
	expected := []byte{0, 0, 0, 15, '/', 'p', 'a', 't', 'h', '/', 't', 'o', '/', 'f', 'i', 'l', 'e'}

	encoded, err := ftc.Encode(&ft)
	if err != nil {
		t.Errorf("Encode() error = %v, wantErr %v", err, nil)
	}
	if !bytes.Equal(encoded, expected) {
		t.Errorf("Encode() got = %v, want %v", encoded, expected)
	}
}

func TestFileTransferCodec_Decode_ShouldReturnFileTransfer_WhenGivenBytes(t *testing.T) {
	ftc := codec.FileTransferCodec{}
	data := []byte{0, 0, 0, 15, '/', 'p', 'a', 't', 'h', '/', 't', 'o', '/', 'f', 'i', 'l', 'e'}
	expected := codec.FileTransfer{
		Length:   15,
		FilePath: "/path/to/file",
	}

	decoded, err := ftc.Decode(data)
	if err != nil {
		t.Errorf("Decode() error = %v, wantErr %v", err, nil)
	}
	if !reflect.DeepEqual(decoded, &expected) {
		t.Errorf("Decode() got = %v, want %v", decoded, &expected)
	}
}

func TestFileTransferAckCodec_Encode_ShouldReturnBytes_WhenGivenFileTransferAck(t *testing.T) {
	ftac := codec.FileTransferAckCodec{}
	fta := codec.FileTransferAck{
		FileID:    1,
		FileLen:   10,
		Checksum:  12345,
		BlockSize: 1024,
		ErrorCode: 0,
	}
	buf := new(bytes.Buffer)
	buf.WriteByte(0)
	buf.WriteByte(0)
	buf.WriteByte(0)
	buf.WriteByte(1) // FileID
	binary.Write(buf, binary.BigEndian, fta.FileLen)
	binary.Write(buf, binary.BigEndian, fta.Checksum)
	binary.Write(buf, binary.BigEndian, fta.BlockSize)
	binary.Write(buf, binary.BigEndian, fta.ErrorCode)
	expected := buf.Bytes()

	encoded, err := ftac.Encode(&fta)
	if err != nil {
		t.Errorf("Encode() error = %v, wantErr %v", err, nil)
	}
	if !bytes.Equal(encoded, expected) {
		t.Errorf("Encode() got = %v, want %v", encoded, expected)
	}
}

func TestFileTransferAckCodec_Decode_ShouldReturnFileTransferAck_WhenGivenBytes(t *testing.T) {
	ftac := codec.FileTransferAckCodec{}
	buf := new(bytes.Buffer)
	buf.WriteByte(0)
	buf.WriteByte(0)
	buf.WriteByte(0)
	buf.WriteByte(1)                           // FileID
	buf.Write([]byte{0, 0, 0, 0, 0, 0, 0, 10}) // FileLen
	buf.Write([]byte{0, 0, 48, 57})            // Checksum
	buf.Write([]byte{0, 0, 4, 0})              // BlockSize
	buf.Write([]byte{0, 0, 0, 0})              // ErrorCode
	data := buf.Bytes()
	expected := codec.FileTransferAck{
		FileID:    1,
		FileLen:   10,
		Checksum:  12345,
		BlockSize: 1024,
		ErrorCode: 0,
	}

	decoded, err := ftac.Decode(data)
	if err != nil {
		t.Errorf("Decode() error = %v, wantErr %v", err, nil)
	}
	if !reflect.DeepEqual(decoded, &expected) {
		t.Errorf("Decode() got = %v, want %v", decoded, &expected)
	}
}

func TestTransferCodec_Encode_ShouldReturnBytes_WhenGivenTransfer(t *testing.T) {
	tc := codec.TransferCodec{}
	tr := codec.Transfer{
		FileID: 1,
		Seq:    2,
		Block:  []byte("Hello, world!"),
	}
	buf := new(bytes.Buffer)
	buf.WriteByte(0)
	buf.WriteByte(0)
	buf.WriteByte(0)
	buf.WriteByte(1) // FileID
	buf.WriteByte(0)
	buf.WriteByte(0)
	buf.WriteByte(0)
	buf.WriteByte(2) // Seq
	binary.Write(buf, binary.BigEndian, uint32(len(tr.Block)))
	buf.Write(tr.Block)
	expected := buf.Bytes()

	encoded, err := tc.Encode(&tr)
	if err != nil {
		t.Errorf("Encode() error = %v, wantErr %v", err, nil)
	}
	if !bytes.Equal(encoded, expected) {
		t.Errorf("Encode() got = %v, want %v", encoded, expected)
	}
}

func TestTransferCodec_Decode_ShouldReturnTransfer_WhenGivenBytes(t *testing.T) {
	tc := codec.TransferCodec{}
	buf := new(bytes.Buffer)
	buf.WriteByte(0)
	buf.WriteByte(0)
	buf.WriteByte(0)
	buf.WriteByte(1) // FileID
	buf.WriteByte(0)
	buf.WriteByte(0)
	buf.WriteByte(0)
	buf.WriteByte(2) // Seq
	blockData := []byte("Hello, world!")
	binary.Write(buf, binary.BigEndian, uint32(len(blockData)))
	buf.Write(blockData)
	data := buf.Bytes()
	expected := codec.Transfer{
		FileID: 1,
		Seq:    2,
		Block:  blockData,
	}

	decoded, err := tc.Decode(data)
	if err != nil {
		t.Errorf("Decode() error = %v, wantErr %v", err, nil)
	}
	if !reflect.DeepEqual(decoded, &expected) {
		t.Errorf("Decode() got = %v, want %v", decoded, &expected)
	}
}
