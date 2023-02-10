package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	logLevel = zapcore.InfoLevel
)

func levelEnabler() zap.LevelEnablerFunc {
	return func(lvl zapcore.Level) bool {
		return lvl >= logLevel
	}
}

func allEnabler() zap.LevelEnablerFunc {
	return func(lvl zapcore.Level) bool {
		return true
	}
}

func errorEnabler() zap.LevelEnablerFunc {
	return func(lvl zapcore.Level) bool {
		return lvl >= zapcore.ErrorLevel
	}
}

func infoEnabler() zap.LevelEnablerFunc {
	return func(lvl zapcore.Level) bool {
		return lvl <= zapcore.WarnLevel
	}
}

func debugEnabler() zap.LevelEnablerFunc {
	return func(lvl zapcore.Level) bool {
		return lvl <= zapcore.DebugLevel
	}
}

func andEnabler(funcs ...zap.LevelEnablerFunc) zap.LevelEnablerFunc {
	return func(lvl zapcore.Level) bool {
		for _, f := range funcs {
			if !f(lvl) {
				return false
			}
		}

		return true
	}
}

//noinspection GoUnusedFunction
func orEnabler(funcs ...zap.LevelEnablerFunc) zap.LevelEnablerFunc {
	return func(lvl zapcore.Level) bool {
		for _, f := range funcs {
			if f(lvl) {
				return true
			}
		}

		return false
	}
}

//noinspection GoUnusedFunction
func SetLogLevel(lvl zapcore.Level) {
	logLevel = lvl
}
