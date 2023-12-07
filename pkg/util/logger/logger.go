package logger

import (
	"log"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	nanoId "github.com/matoous/go-nanoid/v2"
)

// A Logger provides fast, leveled, structured logging. All methods are safe
// for concurrent use, along with filter policy to synthesis logging data.
type Logger struct {
	*zap.Logger
	level         Level
	correlationId string
}

// NewLog create a new logger instance with default Zap logging production config
// and a logging scope based on the given name parameter. Custom logging option
// enables filter policy, correlationID and other configuration for logger.
// Logging is enabled at Info Level and above.
//
// For further logging function. please refer to: https://pkg.go.dev/go.uber.org/zap
//
// Example:
// Create a new logger with name "logger-service", and filter policy
// to redact value from "email" fields.
// logger.NewLog("logger-service", logger.WithFilters(filter.Field("email")))
func NewLog(name string, options ...Option) *Logger {
	result := &Logger{
		level:         LevelInfo,
		correlationId: nanoId.Must(),
	}
	for _, opt := range options {
		opt(result)
	}

	result.Logger = result.newZapLogger(name)

	return result
}

func (l *Logger) newZapLogger(name string) *zap.Logger {
	zapConfig := zap.NewProductionConfig()
	zapConfig.Level = zap.NewAtomicLevelAt(l.level.ToZapLevel())
	zapConfig.EncoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder
	zapConfig.EncoderConfig.EncodeDuration = zapcore.StringDurationEncoder

	zapLogger, err := zapConfig.Build()
	if err != nil {
		log.Println(err)
		return nil
	}

	defer func() {
		_ = zapLogger.Sync()
	}()

	zapLogger = zapLogger.Named(name)
	if l.correlationId != "" {
		zapLogger = zapLogger.With(zap.String("cid", l.correlationId))
	}

	return zapLogger
}

func (l *Logger) clone() *Logger {
	clonedLogger := *l
	return &clonedLogger
}

func (l *Logger) WithErr(err error) *Logger {
	if err == nil {
		return l
	}

	cloned := l.clone()
	cloned.Logger = cloned.With(zap.String("error", err.Error()))

	return cloned
}

func (l *Logger) WithFields(fields map[string]any) *Logger {
	cloned := l.clone()
	logFields := make([]zap.Field, 0)

	for key, value := range fields {
		logFields = append(logFields, zap.Any(key, value))
	}

	cloned.Logger = cloned.With(logFields...)

	return cloned
}

func (l *Logger) Cid() string {
	return l.correlationId
}
