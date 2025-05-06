package logger

import (
	"log/slog"
	"os"
)

// グローバルで使用できるLogger
var Logger *slog.Logger

// init
func init() {
	// json形式でログを出力する
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo, // Infoレベル以上を出力
	})
	Logger = slog.New(handler)
}

// 情報ログを出力
func Info(msg string, keysAndValues ...interface{}) {
	Logger.Info(msg, keysAndValues...)
}

// errorログを出力
func Error(msg string, keysAndValues ...interface{}) {
	Logger.Error(msg, keysAndValues...)
}
