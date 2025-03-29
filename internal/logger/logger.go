package logger

import "go.uber.org/zap"

type Logger struct {
	*zap.SugaredLogger
}

func NewLogger() Logger {
	zl := zap.Must(zap.NewProduction()).Sugar()
	return Logger{zl}
}
