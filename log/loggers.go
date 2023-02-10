package log

import (
	"github.com/alexdogonin/zapsentry"
	"github.com/getsentry/sentry-go"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"rds-enhanced-exporter/config"
)

func infoLogger() *zap.SugaredLogger {
	core := zapcore.NewTee(append(consoleCores(), fileCores(config.GetConfig().Global.Log+"_"+config.GetConfig().Global.Port)...)...)
	logger := zap.New(core, zap.AddCaller()).Named("rds-enhanced-exporter.info")

	return logger.Sugar()
}

func errorLogger() *zap.SugaredLogger {
	core := zapcore.NewTee(append(consoleCores(), fileCores(config.GetConfig().Global.Log+"_"+config.GetConfig().Global.Port)...)...)
	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(errorEnabler())).Named("rds-enhanced-exporter.error")
	reqBodyCtxValName := "body"

	sentryClient, _ := sentry.NewClient(sentry.ClientOptions{
		Dsn:              config.GetConfig().Global.SentryDSN,
		AttachStacktrace: true,
	})

	logger = logger.WithOptions(zap.Development())
	logger = logger.
		WithOptions(
			zap.WrapCore(
				zapsentry.NewWrapper(
					sentryClient,
					zapsentry.WithRequest("request", &reqBodyCtxValName),
					zapsentry.WithSecretHeaders("Authorization-Token"),
				),
			),
		)

	return logger.Sugar()
}
