package message

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
)

type MsgType byte

const FormatVersion = 0

const (
	// TODO: Should I make the zero value be UnknownMsgType or something?
	ErrorMsgType MsgType = iota
	EchoMsgType
)

const MsgHeaderLength = 5

type Header struct {
	Version byte
	Length  int
}

func ParseHeader(data []byte) (Header, bool) {
	if len(data) < MsgHeaderLength {
		return Header{}, false
	}

	return Header{
		Version: data[0],
		Length:  int(binary.BigEndian.Uint32(data[1:5])),
	}, true
}

type Msg interface {
	EncodeMsg() []byte
}

func Send(w io.Writer, m Msg) error {
	_, err := w.Write(m.EncodeMsg())
	return err
}

func Parse(data []byte) (Msg, error) {
	log.Printf("Parsing message from data: %v", data)
	t := data[0]
	switch MsgType(t) {
	case ErrorMsgType:
		return ParseErrorMsg(data[1:]), nil
	case EchoMsgType:
		return ParseEchoMsg(data[1:]), nil
	}
	return nil, fmt.Errorf("Invalid message type: %d", t)
}

func Split(data []byte, _ bool) (advance int, token []byte, err error) {
	log.Println("message.Split")

	header, ok := ParseHeader(data)
	if !ok {
		log.Println("Header not recieved. Waiting for more data")
		// Full header hasn't been recieved.
		// Wait for more data.
		return 0, nil, nil
	}

	if len(data) < int(header.Length) {
		log.Printf("Full message not recieved (have %d bytes, need %d). Waiting for more data", len(data), header.Length)
		// Length of given data is less than length of expected data.
		// The full message hasn't been recieved.
		// Wait for more data.
		return 0, nil, nil
	}

	advance = header.Length
	token = data[:header.Length]

	log.Printf("Got message (return %d, %s, %v)", advance, token, err)

	return
}

type ErrorMsg struct {
	Err string
}

func ParseErrorMsg(data []byte) ErrorMsg {
	return ErrorMsg{Err: string(data)}
}

func (e ErrorMsg) Error() string {
	return e.Err
}

func (self ErrorMsg) EncodeMsg() []byte {
	// Header (5 bytes)
	// Length is empty for now
	data := []byte{FormatVersion, 0x00, 0x00, 0x00, 0x00}
	data = append(data, byte(ErrorMsgType))
	data = append(data, []byte(self.Err)...)

	// Encode total length as uint32 and put it in the length field (previously empty)
	binary.BigEndian.PutUint32(data[1:5], uint32(len(data)))

	return data
}

type EchoMsg struct {
	Message string
}

func ParseEchoMsg(data []byte) EchoMsg {
	return EchoMsg{Message: string(data)}
}

func (self EchoMsg) EncodeMsg() []byte {
	// Header (5 bytes)
	// Length is empty for now
	data := []byte{FormatVersion, 0x00, 0x00, 0x00, 0x00}
	data = append(data, byte(EchoMsgType))
	data = append(data, []byte(self.Message)...)

	// Encode total length as uint32 and put it in the length field (previously empty)
	binary.BigEndian.PutUint32(data[1:5], uint32(len(data)))

	return data
}
