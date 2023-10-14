package logger

import "fmt"

type ConsoleLogger struct {
	Name string
}

func (l *ConsoleLogger) Error(message string) {
	fmt.Printf("\nSERVER: %s [%s]: %s\n", l.Name, "ERROR", message)
}

func (l *ConsoleLogger) Warning(message string) {
	fmt.Printf("\nSERVER: %s\n [%s]: %s\n", l.Name, "WARNING", message)
}

func (l *ConsoleLogger) Info(message string) {
	fmt.Printf("\nSERVER: %s [%s]: %s\n", l.Name, "INFO", message)
}

func NewConsoleLogger(ServerName string) *ConsoleLogger {
	return &ConsoleLogger{ServerName}
}
