package config

import (
	"fmt"
	"github.com/sirupsen/logrus"
)

type ValidatableConf interface {
	IsValid(logger *logrus.Logger, currentPath string, invalidArgs map[string]string) bool
	GetFields() logrus.Fields
}


func AddInvalidArg(currentPath, fieldName, val string, invalidArgs map[string]string) {
	//Add tag to map of invalid args
	path := concatTag(currentPath, fieldName)
	invalidArgs[path] = fmt.Sprintf("<%v> field is using an invalid value <%v>", fieldName, val)
}

func concatTag(current, tag string) string{
	return fmt.Sprintf("%s.%s", current, tag)
}