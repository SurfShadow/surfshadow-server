package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/SurfShadow/surfshadow-server/internal/infrastructure/config"
)

type Logger struct {
	*zap.SugaredLogger
}

func NewLogger(cfg *config.LoggerConfig) (*Logger, error) {
	var zapLevel zapcore.Level
	if err := zapLevel.UnmarshalText([]byte(cfg.Level)); err != nil {
		return nil, err
	}

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:          "time",
		LevelKey:         "level",
		NameKey:          "logger",
		CallerKey:        "caller",
		MessageKey:       "msg",
		EncodeTime:       zapcore.ISO8601TimeEncoder,
		EncodeLevel:      customColorLevelEncoder,
		EncodeCaller:     zapcore.ShortCallerEncoder,
		ConsoleSeparator: " ",
	}

	var core zapcore.Core
	if cfg.Environment == "prod" {
		core = zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderConfig),
			zapcore.AddSync(os.Stdout),
			zap.NewAtomicLevelAt(zapLevel),
		)
	} else {
		core = zapcore.NewCore(
			zapcore.NewConsoleEncoder(encoderConfig),
			zapcore.AddSync(os.Stdout),
			zap.NewAtomicLevelAt(zapLevel),
		)
	}

	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(0))

	return &Logger{logger.Sugar()}, nil
}

func customColorLevelEncoder(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	switch level {
	case zapcore.DebugLevel:
		enc.AppendString("[\033[36mDEBUG\033[0m]")
	case zapcore.InfoLevel:
		enc.AppendString("[\033[32mINFO\033[0m]")
	case zapcore.WarnLevel:
		enc.AppendString("[\033[33mWARN\033[0m]")
	case zapcore.ErrorLevel:
		enc.AppendString("[\033[31mERROR\033[0m]")
	case zapcore.DPanicLevel, zapcore.PanicLevel, zapcore.FatalLevel:
		enc.AppendString("[\033[35mFATAL\033[0m]")
	default:
		enc.AppendString("[UNKNOWN]")
	}
}

func (l *Logger) Sync() {
	_ = l.SugaredLogger.Sync()
}
