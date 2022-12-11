package logger

import (
	"log"
	"os"
)

var logInfo = log.New(os.Stderr, "[🔧] ", 0)
var logErr = log.New(os.Stderr, "[💥] ", 0)

func Info(format string, v ...any) {
	logInfo.Printf(format, v...)
}

func Err(format string, v ...any) {
	logErr.Printf(format, v...)
}

func Fatal(format string, v ...any) {
	logErr.Fatalf(format, v...)
}

func NewFile(filename string, hash string) {
	Info("✨ New: %s[%s]", filename, hash[:7])
}

func DelFile(filename string, hash string) {
	Info("🧹 Del: %s[%s]", filename, hash[:7])
}
