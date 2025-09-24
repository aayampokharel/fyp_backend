package logger

import (
	"fmt"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	Logger *zap.SugaredLogger
	once   sync.Once
)

func InitLogger() {
	once.Do(func() {
		config := zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = customLevelEncoder
		config.EncoderConfig.TimeKey = "T"
		config.EncoderConfig.CallerKey = "C"
		config.EncoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
		config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

		baseLogger, err := config.Build(zap.AddCaller())
		if err != nil {
			panic(err)
		}
		Logger = baseLogger.Sugar()
	})
}

func customLevelEncoder(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	var color string
	var icon string
	switch level {
	case zapcore.DebugLevel:
		color = "\033[35m" // Magenta
		icon = "🔍"
	case zapcore.InfoLevel:
		color = "\033[32m" // Green
		icon = "🚀"
	case zapcore.WarnLevel:
		color = "\033[33m" // Yellow
		icon = "⚠️"
	case zapcore.ErrorLevel, zapcore.DPanicLevel, zapcore.FatalLevel, zapcore.PanicLevel:
		color = "\033[31m" // Red
		icon = "🚨"
	default:
		color = "\033[37m" // White
		icon = "🐛"
	}
	enc.AppendString(fmt.Sprintf("%s %s%s\033[0m", icon, color, level.CapitalString()))
}
