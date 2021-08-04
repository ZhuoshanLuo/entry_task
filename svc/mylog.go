package svc

import (
	"log"
	"os"
)

var (
	DebugLogger *log.Logger
	InfoLogger  *log.Logger
	ErrorLogger *log.Logger
)

func LoggerInit() {
	debugFile, _ := os.Create("./log/debug.txt")
	infoFile, _ := os.Create("./log/info.txt")
	errorFile, _ := os.Create("./log/error.txt")
	DebugLogger = log.New(debugFile, "[DEBUG] ", log.Ldate|log.Ltime|log.Lshortfile)
	InfoLogger = log.New(infoFile, "[INFO] ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(errorFile, "[ERROR] ", log.Ldate|log.Ltime|log.Lshortfile)
}