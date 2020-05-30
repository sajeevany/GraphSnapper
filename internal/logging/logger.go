package logging

import (
	"github.com/sirupsen/logrus"
)

const LoggerKey = "logger"

func Init() *logrus.Logger {
	return logrus.New()
}
