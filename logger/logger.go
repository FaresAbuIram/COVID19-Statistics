package logger

import (
	"log"
	"os"
)

type LoggerCollection struct {
	ErrorLogger *log.Logger
	InfoLogger  *log.Logger
	DebugLogger *log.Logger
}

func NewLoggerCollection() *LoggerCollection {
	return &LoggerCollection{
		ErrorLogger: log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile),
		InfoLogger:  log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime),
		DebugLogger: log.New(os.Stdout, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile),
	}
}

func (l *LoggerCollection) AddErrorLogger(message string) {
	l.ErrorLogger.Println(message)
}

func (l *LoggerCollection) AddInfoLogger(message string) {
	l.InfoLogger.Println(message)
}

func (l *LoggerCollection) AddDebugLogger(message string) {
	l.DebugLogger.Println(message)
}
