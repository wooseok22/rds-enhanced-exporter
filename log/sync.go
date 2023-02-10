package log

import (
	"go.uber.org/zap/zapcore"
	"log"
	"os"
	"path"
)

var logDir = ""

func stdoutSync() zapcore.WriteSyncer {
	return zapcore.Lock(os.Stdout)
}

func stderrSync() zapcore.WriteSyncer {
	return zapcore.Lock(os.Stderr)
}

func fileSync(name string) zapcore.WriteSyncer {
	prefix := getPrefix()
	filepath := path.Join(prefix, name)

	file, err := os.OpenFile(filepath, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0640)

	if err != nil {
		log.Printf("Failed open the log file - %s", filepath)
		return nil
	}

	return zapcore.Lock(file)
}

func getPrefix() string {
	if len(logDir) != 0 {
		return logDir
	}

	return "./"
}

func SetLogDir(dirname string) {
	logDir = dirname
}

func Flush() error {
	errs := make([]error, 0)
	if err := info.Sync(); err != nil {
		errs = append(errs, err)
	}

	if err := errors.Sync(); err != nil {
		errs = append(errs, err)
	}

	if len(errs) > 0 {
		return errs[0]
	}

	return nil
}
