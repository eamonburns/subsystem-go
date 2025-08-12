package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/eamonburns/subsystem-go/internal/message"
)

const (
	ColorRed   = "\033[31m"
	ColorReset = "\033[0m"
)

func main() {
	{ // Setup logging
		logFile, err := os.OpenFile("./server.log", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o666)
		if err != nil {
			panic(err)
		}
		defer logFile.Close()
		log.SetPrefix("[subsystem-server]")
		log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
		log.SetOutput(logFile)
	}

	log.Println("Starting")

	// Using stdio
	//conn := os.Stdin

	// Using network sockets
	listener, err := net.Listen("tcp", "0.0.0.0:8080")
	if err != nil {
		panic(err)
	}
	defer listener.Close()
	conn, err := listener.Accept()
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	scanner := bufio.NewScanner(conn)
	scanner.Split(message.Split)

	for scanner.Scan() {
		msgBytes := scanner.Bytes()

		m, err := message.Parse(msgBytes[message.MsgHeaderLength:])
		if err != nil {
			text := fmt.Sprintf("Parse Error: %v", err)
			log.Println(text)
			fmt.Println(ColorRed + text + ColorReset)
			continue
		}

		switch m := m.(type) {
		case message.ErrorMsg:
			text := fmt.Sprintf("ErrorMsg: %s", m.Err)
			log.Println(text)
			fmt.Println(ColorRed + text + ColorReset)
		case message.EchoMsg:
			text := fmt.Sprintf("EchoMsg: %s", m.Message)
			log.Println(text)
			fmt.Println(text)
		default:
			text := "Error: Unknown message"
			log.Println(text)
			fmt.Println(ColorRed + text + ColorReset)
		}
	}

	err = scanner.Err()
	if err != nil {
		panic(err)
	}
}
