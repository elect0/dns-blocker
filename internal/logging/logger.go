package logging

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"strings"
	"time"
)

type Level int8

const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal
)

func (l Level) String() string {
	switch l {
	case LevelDebug:
		return "DEBUG"
	case LevelInfo:
		return "INFO"
	case LevelWarn:
	return "WARN"
	case LevelError:
		return "ERROR"
	case LevelFatal:
		return "FATAL"
	}
	return ""
}

type Logger struct {
	out      *log.Logger
	minLevel Level
}

func New(out io.Writer, minLevel Level) *Logger {
	return &Logger{
		out:      log.New(out, "", 0),
		minLevel: minLevel,
	}
}

func StringToLevel(levelStr string) (Level, error) {
	switch strings.ToUpper(levelStr) {
	case "DEBUG":
		return LevelDebug, nil
	case "INFO":
		return LevelInfo, nil
	case "WARN":
		return LevelWarn, nil
	case "ERROR":
		return LevelError, nil
	case "FATAL":
		return LevelFatal, nil
	}
	return LevelDebug, fmt.Errorf("invalid log level: %v", levelStr)
}

func (l *Logger) print(level Level, message string, properties map[string]string) {
	if level < l.minLevel {
		return
	}

	logEntry := struct {
		Timestamp  string            `json:"ts"`
		Level      string            `json:"level"`
		Message    string            `json:"msg"`
		Properties map[string]string `json:"properties"`
		Caller     string            `json:"caller"`
	}{
		Timestamp:  time.Now().UTC().Format(time.RFC3339),
		Level:      level.String(),
		Message:    message,
		Properties: properties,
	}

	_, file, line, ok := runtime.Caller(2)
	if ok {
		logEntry.Caller = fmt.Sprintf("%s:%d", file, line)
	}

	lineBytes, err := json.Marshal(logEntry)
	if err != nil {
		l.out.Printf("{\"level\":\"ERROR\", \"msg\":\"failed to marshal log message: %v\"}", err)
		return
	}

	l.out.Println(string(lineBytes))

	if level == LevelFatal {
		os.Exit(1)
	}
}

func (l *Logger) Debug(message string, properties map[string]string) {
	l.print(LevelDebug, message, properties)
}

func (l *Logger) Info(message string, properties map[string]string) {
	l.print(LevelInfo, message, properties)
}

func (l *Logger) Warn(message string, properties map[string]string) {
	l.print(LevelWarn, message, properties)
}

func (l *Logger) Error(message string, err error, properties map[string]string) {
	if properties == nil {
		properties = make(map[string]string)
	}

	properties["error"] = err.Error()
	l.print(LevelError, message, properties)
}

func (l *Logger) Fatal(message string, err error, properties map[string]string) {
	if properties == nil {
		properties = make(map[string]string)
	}

	properties["error"] = err.Error()
	l.print(LevelFatal, message, properties)
}
