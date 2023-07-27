package initializers

import (
	"log"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.SugaredLogger

const (
	InfoLevel  = zapcore.InfoLevel  // InfoLevel logs are Debug priority logs.
	WarnLevel  = zapcore.WarnLevel  // WarnLevel logs are more important than Info, but don't need individual human review.
	ErrorLevel = zapcore.ErrorLevel // ErrorLevel logs are high-priority.
	FatalLevel = zapcore.FatalLevel // FatalLevel logs a message, then calls os.Exit(1).
)

const (
	InfoLogFilePath  = "logs/info.log"
	WarnLogFilePath  = "logs/warn.log"
	ErrorLogFilePath = "logs/error.log"
	FatalLogFilePath = "logs/fatal.log"
)

var logFiles struct {
	InfoLogFile  *os.File
	WarnLogFile  *os.File
	ErrorLogFile *os.File
	FatalLogFile *os.File
}

func AddLogger() {
	openLogFiles()

	infoCore := newCore(logFiles.InfoLogFile, InfoLevel)
	warnCore := newCore(logFiles.WarnLogFile, WarnLevel)
	errorCore := newCore(logFiles.ErrorLogFile, ErrorLevel)
	fatalCore := newCore(logFiles.FatalLogFile, FatalLevel)

	Logger = zap.New(zapcore.NewTee(infoCore, warnCore, errorCore, fatalCore)).Sugar()
}

func newCore(LogFile *os.File, LoggerLevel zapcore.Level) zapcore.Core {
	encoderCfg := zap.NewProductionEncoderConfig()
	fileCore := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderCfg),
		zapcore.AddSync(LogFile),
		LoggerLevel,
	)
	return fileCore
}

func openLogFiles() {
	logFiles.InfoLogFile = openFile(InfoLogFilePath)
	logFiles.WarnLogFile = openFile(WarnLogFilePath)
	logFiles.ErrorLogFile = openFile(ErrorLogFilePath)
	logFiles.FatalLogFile = openFile(FatalLogFilePath)
}

func openFile(LogFilePath string) *os.File {
	logFile, err := os.OpenFile(LogFilePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Printf("Failed to open log file: " + err.Error())
	}
	return logFile
}

func LoggerCleanUp() {
	logFiles.InfoLogFile.Close()
	logFiles.WarnLogFile.Close()
	logFiles.ErrorLogFile.Close()
	logFiles.FatalLogFile.Close()

	Logger.Sync()
}
