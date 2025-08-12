package main

import (
	"bufio"
	"log"
	"net"
	"os"
	"strings"

	"github.com/eamonburns/subsystem-go/internal/message"
)

func main() {
	{ // Setup logging
		logFile, err := os.OpenFile("./client.log", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o666)
		if err != nil {
			panic(err)
		}
		defer logFile.Close()
		log.SetPrefix("[subsystem-client]")
		log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
		log.SetOutput(logFile)
	}

	log.Println("Connecting")

	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		panic(err)
	}

	log.Printf("Connected. Sending messages")

	message.Send(conn, message.EchoMsg{
		Message: "This is an echo",
	})
	message.Send(conn, message.ErrorMsg{
		Err: "This is an error",
	})
	message.Send(conn, message.EchoMsg{
		Message: "Echo 2",
	})

	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		line := input.Text()
		var m message.Msg
		if strings.HasPrefix(line, "e ") {
			m = message.ErrorMsg{Err: line[2:]}
		} else if strings.HasPrefix(line, "m ") {
			m = message.EchoMsg{Message: line[2:]}
		} else {
			// TODO: handle it better
			panic("Unknown type")
		}

		if err := message.Send(conn, m); err != nil {
			panic(err)
		}
	}
}
