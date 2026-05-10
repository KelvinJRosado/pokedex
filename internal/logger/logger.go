package logger

import (
	"log/slog"
	"os"
	"path/filepath"
	"time"
)

// Holds both a stdio logger + file logger
type CustomLogger struct {
	consoleWriter *slog.Logger
	logFileWriter *slog.Logger
	file          *os.File
}

func NewLogger() (*CustomLogger, error) {

	// Get file location from env var, or use default
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	logPath := filepath.Join(cwd, "logs")

	// Ensure log directory exists
	err = os.MkdirAll(logPath, 0755)
	if err != nil {
		return nil, err
	}

	// Generate log file name
	logFileName := time.Now().Format("2006-01-02")
	logFilePath := filepath.Join(logPath, logFileName+".log")

	// Open log file
	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}

	// Build handlers
	loggerOptions := &slog.HandlerOptions{Level: slog.LevelDebug}
	textHandler := slog.NewTextHandler(os.Stdout, loggerOptions) // text to console
	jsonHandler := slog.NewJSONHandler(logFile, loggerOptions)   // JSON to log file

	// Build loggers
	textLogger := slog.New(textHandler)
	jsonLogger := slog.New(jsonHandler)

	logger := CustomLogger{consoleWriter: textLogger, logFileWriter: jsonLogger, file: logFile}
	return &logger, nil
}

// Close the log file when done
func (cl *CustomLogger) Close() error {
	return cl.file.Close()
}

func (cl *CustomLogger) Debug(msg string, args ...any) {
	cl.consoleWriter.Debug(msg, args...)
	cl.logFileWriter.Debug(msg, args...)
}

func (cl *CustomLogger) Info(msg string, args ...any) {
	cl.consoleWriter.Info(msg, args...)
	cl.logFileWriter.Info(msg, args...)
}

func (cl *CustomLogger) Warn(msg string, args ...any) {
	cl.consoleWriter.Warn(msg, args...)
	cl.logFileWriter.Warn(msg, args...)
}

func (cl *CustomLogger) Error(msg string, args ...any) {
	cl.consoleWriter.Error(msg, args...)
	cl.logFileWriter.Error(msg, args...)
}
