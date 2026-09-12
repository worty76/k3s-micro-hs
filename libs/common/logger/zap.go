package logger

import (
	"os"

	"github.com/worty76/k3s-micro-hs/libs/common/env"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Zap struct {
	zapLogger *zap.Logger
}

func NewZap(environment env.Environment) *Zap {
	var config zapcore.EncoderConfig

	if environment == env.Prod {
		config = zap.NewProductionEncoderConfig()
		config.EncodeTime = zapcore.ISO8601TimeEncoder
	} else {
		config = zap.NewDevelopmentEncoderConfig()
		config.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	// Microservice Production defaults
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(config), // JSON is native to log aggregators like ELK/Grafana Loki
		zapcore.Lock(os.Stdout),
		zap.NewAtomicLevelAt(zapcore.InfoLevel), // Default microservice log level
	)

	// Wrap core with options optimal for professional environments
	logger := zap.New(core,
		zap.AddCaller(),                       // Shows file line numbers
		zap.AddCallerSkip(1),                  // Important: Skips the adapter abstraction wrapper layer
		zap.AddStacktrace(zapcore.ErrorLevel), // Automatic stack traces on errors
	)

	return &Zap{zapLogger: logger}
}

func (z *Zap) toZapFields(fields []Field) []zap.Field {
	zapFields := make([]zap.Field, len(fields))
	for i, f := range fields {
		// zap.Any still use reflection but it is the most flexible way to log arbitrary types
		zapFields[i] = zap.Any(f.Key, f.Value)
	}
	return zapFields
}

func (z *Zap) Info(msg string, fields ...Field) {
	z.zapLogger.Info(msg, z.toZapFields(fields)...)
}

func (z *Zap) Error(msg string, fields ...Field) {
	z.zapLogger.Error(msg, z.toZapFields(fields)...)
}

func (z *Zap) Debug(msg string, fields ...Field) {
	z.zapLogger.Debug(msg, z.toZapFields(fields)...)
}

func (z *Zap) Warn(msg string, fields ...Field) {
	z.zapLogger.Warn(msg, z.toZapFields(fields)...)
}

func (z *Zap) Fatal(msg string, fields ...Field) {
	z.zapLogger.Fatal(msg, z.toZapFields(fields)...)
}
