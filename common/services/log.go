package services

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
)

type LogService struct {
	file     *os.File
	mu       sync.Mutex
	logPath  string
	logLevel LogLevel
}

var (
	logService *LogService
	once       sync.Once
)

func NewLogService() *LogService {
	once.Do(func() {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			panic(fmt.Sprintf("Failed to get user home directory: %v", err))
		}

		logDir := filepath.Join(homeDir, ".fg", "logs")
		if err := os.MkdirAll(logDir, 0755); err != nil {
			panic(fmt.Sprintf("Failed to create log directory: %v", err))
		}

		logPath := filepath.Join(logDir, "app.log")
		file, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			panic(fmt.Sprintf("Failed to open log file: %v", err))
		}

		logService = &LogService{
			file:     file,
			logPath:  logPath,
			logLevel: INFO,
		}
	})
	return logService
}

func (ls *LogService) SetLogLevel(level LogLevel) {
	ls.mu.Lock()
	defer ls.mu.Unlock()
	ls.logLevel = level
}

func (ls *LogService) log(level LogLevel, format string, args ...interface{}) {
	if level < ls.logLevel {
		return
	}

	ls.mu.Lock()
	defer ls.mu.Unlock()

	levelStr := "UNKNOWN"
	switch level {
	case DEBUG:
		levelStr = "DEBUG"
	case INFO:
		levelStr = "INFO"
	case WARN:
		levelStr = "WARN"
	case ERROR:
		levelStr = "ERROR"
	}

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	message := fmt.Sprintf(format, args...)
	logEntry := fmt.Sprintf("[%s] [%s] %s\n", timestamp, levelStr, message)

	if _, err := ls.file.WriteString(logEntry); err != nil {
		fmt.Printf("Failed to write log: %v\n", err)
	}
}

func (ls *LogService) Debug(format string, args ...interface{}) {
	ls.log(DEBUG, format, args...)
}

func (ls *LogService) Info(format string, args ...interface{}) {
	ls.log(INFO, format, args...)
}

func (ls *LogService) Warn(format string, args ...interface{}) {
	ls.log(WARN, format, args...)
}

func (ls *LogService) Error(format string, args ...interface{}) {
	ls.log(ERROR, format, args...)
}

func (ls *LogService) Close() error {
	ls.mu.Lock()
	defer ls.mu.Unlock()
	return ls.file.Close()
}
