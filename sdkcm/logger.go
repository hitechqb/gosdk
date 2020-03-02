package sdkcm

import (
	"github.com/sirupsen/logrus"
	"sync"
)

var (
	loggerOnce sync.Once
	logger     *logrus.Logger
)

func NewLogger() *logrus.Logger {
	logger = logrus.New()
	//logger.Formatter = &logrus.TextFormatter{
	//	ForceColors:               true,
	//	DisableColors:             false,
	//	EnvironmentOverrideColors: true,
	//	DisableTimestamp:          true,
	//}
	//...
	return logger
}

func GetLogger() *logrus.Logger {
	loggerOnce.Do(func() {
		logger = NewLogger()
	})

	return logger
}
