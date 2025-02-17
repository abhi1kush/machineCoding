package logger

import (
	"io"
	"log"
	"os"

	"gopkg.in/natefinch/lumberjack.v2"
)

var Logger *log.Logger

func InitLogger(logFile string, maxSize, maxBackups, maxAge int, compress bool) {
	lumberjackLogger := &lumberjack.Logger{
		Filename:   logFile,
		MaxSize:    maxSize,    // megabytes
		MaxBackups: maxBackups, // number of backups
		MaxAge:     maxAge,     // days
		Compress:   compress,   // compress rotated files
	}

	multiWriter := io.MultiWriter(os.Stdout, lumberjackLogger)
	Logger = log.New(multiWriter, "", log.LstdFlags|log.Lshortfile)
}
