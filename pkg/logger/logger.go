package logger

import (
	"fmt"
	"io"
	"log"
	"os"

	"forum/internal/conf"
)

type Logger struct {
	*log.Logger
}

func Setup(config *conf.Config) (*Logger, error) {
	file, err := os.OpenFile(config.LogFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o666)
	if err != nil {
		return nil, fmt.Errorf("open logfile - %v", err)
	}

	multiWriter := io.MultiWriter(file, os.Stdout)

	l := &Logger{
		Logger: log.New(multiWriter, "", log.Ldate|log.Ltime|log.Lshortfile),
	}

	return l, nil
}

func (l *Logger) Fatal(v ...interface{}) {
	l.SetPrefix("FATAL: ")
	l.Logger.Fatal(v...)
}

func (l *Logger) Info(v ...interface{}) {
	l.SetPrefix("INFO: ")
	l.Logger.Println(v...)
}

func (l *Logger) Error(v ...interface{}) {
	l.SetPrefix("ERROR: ")
	l.Logger.Println(v...)
}
