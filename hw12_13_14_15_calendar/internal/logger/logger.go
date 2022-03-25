package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var pathToLogFile string

type Logger struct {
	SugarLogger *zap.SugaredLogger
}

func New(logFile, level string) *Logger {
	pathToLogFile = logFile
	writeSyncer := getLogWriter()
	encoder := getEncoder()
	core := zapcore.NewCore(encoder, writeSyncer, getLogLevel(level))

	logger := zap.New(core)

	return &Logger{
		SugarLogger: logger.Sugar(),
	}
}

func (l Logger) Info(msg string) {
	l.SugarLogger.Infof(msg)
}

func (l Logger) Error(msg string) {
	l.SugarLogger.Errorf(msg)
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	return zapcore.NewConsoleEncoder(encoderConfig)
}

func getLogWriter() zapcore.WriteSyncer {
	file, err := os.OpenFile(pathToLogFile, os.O_CREATE|os.O_RDWR, 0o666)
	if err != nil {
		panic("can't create logger: " + err.Error())
	}
	return zapcore.AddSync(file)
}

func getLogLevel(l string) zapcore.Level {
	switch l {
	case "INFO":
		return zapcore.InfoLevel
	case "DEBUG":
		return zapcore.DebugLevel
	case "ERROR":
		return zapcore.ErrorLevel
	default:
		return zapcore.InfoLevel
	}
}
