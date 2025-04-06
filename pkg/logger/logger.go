package logger

import (
	"io"
	"os"
	"path/filepath"
	"runtime"

	"github.com/sirupsen/logrus"
)

// Уровни логирования
const (
	DebugLevel = "debug"
	InfoLevel  = "info"
	WarnLevel  = "warn"
	ErrorLevel = "error"
	FatalLevel = "fatal"
)

// Logger - интерфейс логгера
type Logger interface {
	Debug(args ...interface{})
	Debugf(format string, args ...interface{})
	Info(args ...interface{})
	Infof(format string, args ...interface{})
	Warn(args ...interface{})
	Warnf(format string, args ...interface{})
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
	Fatal(args ...interface{})
	Fatalf(format string, args ...interface{})
}

// LogrusLogger - реализация логгера на основе logrus
type LogrusLogger struct {
	logger *logrus.Logger
}

// NewLogger - создание нового логгера
func NewLogger(level string, output io.Writer) Logger {
	logger := logrus.New()

	// Настройка уровня логирования
	switch level {
	case DebugLevel:
		logger.SetLevel(logrus.DebugLevel)
	case InfoLevel:
		logger.SetLevel(logrus.InfoLevel)
	case WarnLevel:
		logger.SetLevel(logrus.WarnLevel)
	case ErrorLevel:
		logger.SetLevel(logrus.ErrorLevel)
	case FatalLevel:
		logger.SetLevel(logrus.FatalLevel)
	default:
		logger.SetLevel(logrus.InfoLevel)
	}

	// Настройка форматирования
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})

	// Настройка вывода
	if output != nil {
		logger.SetOutput(output)
	} else {
		logger.SetOutput(os.Stdout)
	}

	return &LogrusLogger{
		logger: logger,
	}
}

// NewFileLogger - создание логгера с выводом в файл
func NewFileLogger(level, filePath string) (Logger, error) {
	// Создаем директорию для логов, если она не существует
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}

	// Открываем файл для записи логов
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	return NewLogger(level, file), nil
}

// Debug - логирование отладочных сообщений
func (l *LogrusLogger) Debug(args ...interface{}) {
	l.logger.WithFields(getCallerFields()).Debug(args...)
}

// Debugf - форматированное логирование отладочных сообщений
func (l *LogrusLogger) Debugf(format string, args ...interface{}) {
	l.logger.WithFields(getCallerFields()).Debugf(format, args...)
}

// Info - логирование информационных сообщений
func (l *LogrusLogger) Info(args ...interface{}) {
	l.logger.WithFields(getCallerFields()).Info(args...)
}

// Infof - форматированное логирование информационных сообщений
func (l *LogrusLogger) Infof(format string, args ...interface{}) {
	l.logger.WithFields(getCallerFields()).Infof(format, args...)
}

// Warn - логирование предупреждений
func (l *LogrusLogger) Warn(args ...interface{}) {
	l.logger.WithFields(getCallerFields()).Warn(args...)
}

// Warnf - форматированное логирование предупреждений
func (l *LogrusLogger) Warnf(format string, args ...interface{}) {
	l.logger.WithFields(getCallerFields()).Warnf(format, args...)
}

// Error - логирование ошибок
func (l *LogrusLogger) Error(args ...interface{}) {
	l.logger.WithFields(getCallerFields()).Error(args...)
}

// Errorf - форматированное логирование ошибок
func (l *LogrusLogger) Errorf(format string, args ...interface{}) {
	l.logger.WithFields(getCallerFields()).Errorf(format, args...)
}

// Fatal - логирование критических ошибок с завершением программы
func (l *LogrusLogger) Fatal(args ...interface{}) {
	l.logger.WithFields(getCallerFields()).Fatal(args...)
}

// Fatalf - форматированное логирование критических ошибок с завершением программы
func (l *LogrusLogger) Fatalf(format string, args ...interface{}) {
	l.logger.WithFields(getCallerFields()).Fatalf(format, args...)
}

// getCallerFields - получение информации о вызывающей функции
func getCallerFields() logrus.Fields {
	_, file, line, ok := runtime.Caller(2)
	if !ok {
		file = "???"
		line = 0
	}

	return logrus.Fields{
		"file": filepath.Base(file),
		"line": line,
	}
}
