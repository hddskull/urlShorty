package utils

import "go.uber.org/zap"

var SugaredLogger = newSugaredLogger()

func newSugaredLogger() zap.SugaredLogger {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}

	defer logger.Sync()

	return *logger.Sugar()
}
