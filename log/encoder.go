package log

import (
	"go.uber.org/zap/zapcore"
)

func encoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}

func consoleEncoder() zapcore.Encoder {
	conf := encoderConfig()

	return zapcore.NewConsoleEncoder(conf)
}

func accessLogEncoder() zapcore.Encoder {
	conf := encoderConfig()

	return newAccessEncoder(conf)
}

func jsonEncoder() zapcore.Encoder {
	conf := encoderConfig()
	conf.EncodeTime = zapcore.EpochNanosTimeEncoder
	conf.EncodeDuration = zapcore.NanosDurationEncoder

	return zapcore.NewJSONEncoder(conf)
}
