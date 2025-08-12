package logging

import (
	"log"
	"os"
)

func Setup(filename string) *log.Logger {
	logfile, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o666)
	if err != nil {
		panic(err)
	}

	return log.New(logfile, "[subsystem-go]", log.Ldate|log.Ltime|log.Lshortfile)
}
