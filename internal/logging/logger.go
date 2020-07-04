package logging

import (
	"github.com/sirupsen/logrus"
)

const LoggerKey = "logger"

func Init() *logrus.Logger {
	return logrus.New()
}

//Returns values as a redacted string if non empty
func RedactNonEmpty(val string) string {

	if val != "" {
		return "*****"
	}

	return val
}
