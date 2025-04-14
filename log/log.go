package log

import (
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.Logger

func LogInit() {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.MessageKey = "message"
	encoderConfig.TimeKey = "timestamp"
	encoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
	}
	encoderConfig.EncodeDuration = func(d time.Duration, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendFloat64(float64(d) / float64(time.Millisecond))
	}

	cfg := zap.Config{
		Level:             zap.NewAtomicLevelAt(zapcore.DebugLevel),
		DisableStacktrace: false,
		DisableCaller:     false,
		OutputPaths:       []string{"stdout"},
		Encoding:          "console",
		EncoderConfig:     encoderConfig,
	}

	var err error
	logger, err = cfg.Build(zap.AddStacktrace(zapcore.PanicLevel), zap.AddCallerSkip(1))
	if err != nil {
		panic(err)
	}
}

func Infof(template string, args ...interface{}) {
	logger.Sugar().Infof(template, args...)
}

func Errorf(template string, args ...interface{}) {
	logger.Sugar().Errorf(template, args...)
}

func Info(template string) {
	logger.Sugar().Info(template)
}

func Error(template string) {
	logger.Sugar().Error(template)
}
