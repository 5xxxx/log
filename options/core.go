package options

import (
	"os"
	"time"

	"github.com/mattn/go-colorable"
	"go.elastic.co/ecszap"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func ZapCore(option Options) zapcore.Core {
	switch option.LogType {
	case Console:
		return consoleCore(option)
	case Elastic:
		return elasticCore(option)
	case File:
		return fileCore(option)
	}
	return nil
}

func fileCore(option Options) zapcore.Core {
	cfgConsole := zapcore.EncoderConfig{
		EncodeTime: func(t time.Time, encoder zapcore.PrimitiveArrayEncoder) {
			encoder.AppendString("2006-01-02 15:04:05")
		},
		TimeKey:        "time",
		LevelKey:       "Level",
		NameKey:        "logger",
		CallerKey:      "linenum",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeName:     zapcore.FullNameEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
	atomicLevel := zap.NewAtomicLevel()
	atomicLevel.SetLevel(option.Level)

	return zapcore.NewCore(zapcore.NewConsoleEncoder(cfgConsole),
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout),
			zapcore.AddSync(lumberjackLogger(option))),
		atomicLevel)
}

func consoleCore(option Options) zapcore.Core {
	aa := zap.NewDevelopmentEncoderConfig()
	aa.EncodeLevel = zapcore.CapitalColorLevelEncoder
	aa.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString("2006-01-02 15:04:05")
	}

	return zapcore.NewCore(
		zapcore.NewConsoleEncoder(aa),
		zapcore.AddSync(colorable.NewColorableStdout()),
		option.Level,
	)
}

func elasticCore(option Options) zapcore.Core {
	if option.out == nil {
		panic("writer can't be nil")
	}

	var topicErrors = zapcore.AddSync(option.out) //nolint:gofumpt
	encoderConfig := ecszap.NewDefaultEncoderConfig()
	return ecszap.NewCore(encoderConfig, topicErrors, zap.LevelEnablerFunc(func(level zapcore.Level) bool {
		return option.Level <= level
	}))
}
