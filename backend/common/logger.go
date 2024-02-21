package common

import (
	"os"

	"github.com/wailsapp/wails/v2/pkg/logger"

	"github.com/sirupsen/logrus"
)

type WailsLoggerWrapper struct {
	internalLogger *logrus.Logger
}

// Print Adhering to the Wails Logger interface
func (w *WailsLoggerWrapper) Print(message string) {
	w.internalLogger.Print(message)
}

func (w *WailsLoggerWrapper) Trace(message string) {
	w.internalLogger.Trace(message)
}

func (w *WailsLoggerWrapper) Debug(message string) {
	w.internalLogger.Debug(message)
}

func (w *WailsLoggerWrapper) Info(message string) {
	w.internalLogger.Info(message)
}

func (w *WailsLoggerWrapper) Warning(message string) {
	w.internalLogger.Warn(message)
}

func (w *WailsLoggerWrapper) Error(message string) {
	w.internalLogger.Error(message)
}

func (w *WailsLoggerWrapper) Fatal(message string) {
	w.internalLogger.Fatal(message)
}

// NewLogger Initialize our logger with WailsLoggerWrapper
func NewLogger(logFile string) logger.Logger {
	f, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0664)
	if err != nil {
		panic("创建日志文件失败: " + err.Error())
	}

	logrus.SetOutput(f)
	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors:     true,
		TimestampFormat: "2006-01-02 15:04:05",
		FullTimestamp:   true,
	})

	customLogger := &WailsLoggerWrapper{
		internalLogger: logrus.StandardLogger(),
	}

	return customLogger
}
