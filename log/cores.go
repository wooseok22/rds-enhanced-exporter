package log

import (
	"go.uber.org/zap/zapcore"
)

func consoleCores() []zapcore.Core {
	cores := make([]zapcore.Core, 0)

	cores = append(cores, zapcore.NewCore(consoleEncoder(), stdoutSync(), andEnabler(levelEnabler(), infoEnabler())))
	cores = append(cores, zapcore.NewCore(consoleEncoder(), stderrSync(), andEnabler(levelEnabler(), errorEnabler())))
	cores = append(cores, zapcore.NewCore(consoleEncoder(), stderrSync(), andEnabler(levelEnabler(), debugEnabler())))

	return cores
}

func fileCores(filename string) []zapcore.Core {
	cores := make([]zapcore.Core, 0)

	sync := fileSync(filename)

	if sync == nil {
		return cores
	}

	cores = append(cores, zapcore.NewCore(jsonEncoder(), sync, levelEnabler()))

	return cores
}
