package log

import (
	"io"
	"log"
	"os"
)

var (
	Debug *log.Logger
	Info  *log.Logger
	Warn  *log.Logger
	Error *log.Logger
)

func Init() {
	logLevel := os.Getenv("LOG_LEVEL")

	Debug = log.New(os.Stdout, "[DEBUG] ", log.Ldate|log.Ltime|log.Lshortfile)
	Info  = log.New(os.Stdout, "[INFO] ", log.Ldate|log.Ltime|log.Lshortfile)
	Warn  = log.New(os.Stdout, "[WARN] ", log.Ldate|log.Ltime|log.Lshortfile)
	Error = log.New(os.Stderr, "[ERROR] ", log.Ldate|log.Ltime|log.Lshortfile)

	// Set Debug logger to discard if LOG_LEVEL != debug
	if logLevel != "debug" && logLevel != "DEBUG" {
		Debug.SetOutput(io.Discard)
	}
}

func init() {
	Init()
}
