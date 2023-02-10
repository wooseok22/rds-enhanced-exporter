package log

import "go.uber.org/zap"

var (
	info   *zap.SugaredLogger
	errors *zap.SugaredLogger
)

func Initialize() {
	Create()
}

func Create() {
	info = infoLogger()
	errors = errorLogger()
}

func INFO() *zap.SugaredLogger {
	return info
}

func ERROR() *zap.SugaredLogger {
	return errors
}
