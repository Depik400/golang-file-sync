package logger

type ILogger interface {
	Error(string)
	Warning(string)
	Info(string)
}
