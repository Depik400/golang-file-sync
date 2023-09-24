package logger

import (
	"fmt"
	"log"
	"os"
	"time"
)

type FileLogger struct {
	FileHandle *os.File
	Name       string
}

type FileLoggerConfig struct {
	Path     string
	Rotating bool
	Name     string
}

func (l *FileLogger) Error(message string) {
	l.FileHandle.WriteString(fmt.Sprintf("SERVER: %s [%s]: %s\n", l.Name, "ERROR", message))
}

func (l *FileLogger) Warning(message string) {
	l.FileHandle.WriteString(fmt.Sprintf("SERVER: %s [%s]: %s\n", l.Name, "WARNING", message))
}

func (l *FileLogger) Info(message string) {
	l.FileHandle.WriteString(fmt.Sprintf("SERVER: %s [%s]: %s\n", l.Name, "INFO", message))
}

func NewFileLogger(config *FileLoggerConfig) *FileLogger {
	path := config.Path
	if path[len(path)-1:][0] != '/' {
		path += "/"
	}
	_, err := os.Stat(config.Path)
	if err != nil {
		err := os.Mkdir(config.Path, os.ModePerm)
		if err != nil {
			log.Fatal(err)
		}
	}

	filename := "filelogger"
	if config.Rotating {
		filename = filename + "-" + time.Now().Format("2006-01-02")
	}
	path += filename + ".log"
	_, err = os.Stat(path)
	handle, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.ModeType)
	if err != nil {
		log.Fatal(err)
	}
	return &FileLogger{
		handle,
		config.Name,
	}
}
