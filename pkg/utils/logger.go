package utils

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var LogLevel = zapcore.InfoLevel

func NewLogger() *zap.Logger {
	lowPriority := zap.LevelEnablerFunc(func(lv zapcore.Level) bool {
		return lv >= LogLevel
	})

	jsonEnc := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	stdCore := zapcore.NewCore(jsonEnc, zapcore.Lock(os.Stdout), lowPriority)

	return zap.New(stdCore).WithOptions(zap.AddCaller())
}
