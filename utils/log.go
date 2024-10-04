package utils

import "go.uber.org/zap"

var logger *zap.SugaredLogger

func Errorf(template string, args ...interface{}) {
	logger.Errorf(template, args...)
}
