package util

import (
	"log"
	"os"
)

var (
	ILog *log.Logger
	WLog *log.Logger
	ELog *log.Logger
)

func init () {
	ILog = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	WLog = log.New(os.Stdout, "WARN: ", log.Ldate|log.Ltime|log.Lshortfile)
	ELog = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}
