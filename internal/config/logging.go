package config

import "github.com/sirupsen/logrus"

type Logging struct {
	Level string `json:"level"`
}

func (l Logging) GetFields() logrus.Fields {
	return logrus.Fields{
		"level": l.Level,
	}
}

//IsValid - Returns true/false and a non-empty map of all invalid args. Nested args are set in the form of Parent.Child.SubChild
//Inputs:
//    currentPath - json path defined up and including this attribute. ie conf.Aero
//    invalidArgs - map of invalid arguments (currentPath + field name) mapped to invalid reasons
func (l Logging) IsValid(logger *logrus.Logger, currentPath string, invalidArgs map[string]string) bool {

	isValid := true

	//Check attributes
	if isLoggingLevelInvalid(l.Level) {
		AddInvalidArg(currentPath, "Level", l.Level, invalidArgs)
		isValid = false
	}

	return isValid
}

func isLoggingLevelInvalid(level string) bool {
	_, err := logrus.ParseLevel(level)
	return err != nil
}
