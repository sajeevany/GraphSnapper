package config

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"reflect"
	"strconv"
	"strings"
)

type AerospikeCfg struct {
	Host           string `json:"host"`
	Port           int    `json:"port"`
	Password       string `json:"password"`
	GraphNamespace string `json:"graphNamespace"`
}

func (as AerospikeCfg) getFields() logrus.Fields {
	return logrus.Fields{
		"host":           as.Host,
		"port":           as.Port,
		"password":       as.Password,
		"graphNamespace": as.GraphNamespace,
	}
}

//IsValid - Returns true/false and a non-empty map of all invalid args. Nested args are set in the form of Parent.Child.SubChild
//Inputs:
//    currentPath - json path defined up and including this attribute. ie conf.Aero
func (as AerospikeCfg) IsValid(logger *logrus.Logger, currentPath string) (bool, map[string]string){

	invalidArgs := make(map[string]string, 0)

	//Check attributes
	if as.Host == "" {
		as.addInvalidArg(logger, currentPath, "Host", "Host", as.Host, invalidArgs)
	}

	if as.Port <= 0 || as.Port > 65535{
		as.addInvalidArg(logger, currentPath, "Port", "Port", strconv.Itoa(as.Port), invalidArgs)
	}

	if as.Password == ""{
		as.addInvalidArg(logger, currentPath, "Password", "Password", as.Password, invalidArgs)
	}

	if as.GraphNamespace == ""{
		as.addInvalidArg(logger, currentPath, "GraphNamespace", "GraphNamespace", as.GraphNamespace, invalidArgs)
	}

	return len(invalidArgs) == 0, invalidArgs
}

func (as AerospikeCfg) addInvalidArg(logger *logrus.Logger,  currentPath, fieldName, defaultJsonTagName, val string, invalidArgs map[string]string) {
	//Get json tag as defined by Aerospike cfg
	jsonTag, ok := as.getFirstJsonTagElement(fieldName)
	if !ok {
		logger.WithFields(as.getFields()).Errorf("Unable to get json tag for field <%v>. Defaulting to <%v>", fieldName, defaultJsonTagName)
		jsonTag = defaultJsonTagName
	}

	//Add tag to map of invalid args
	path := concatTag(currentPath, jsonTag)
	invalidArgs[path] = fmt.Sprintf("<%v> field is using an invalid value <%v>", jsonTag, val)
}

func concatTag(current, tag string) string{
	return fmt.Sprintf("%s.%s", current, tag)
}

func (as *AerospikeCfg) getFirstJsonTagElement(fieldName string) (tag string, ok bool) {

	field, ok := reflect.TypeOf(as).Elem().FieldByName(fieldName);
	if !ok{
		return "", false
	}

	//Get json tag
	jTag, ok := field.Tag.Lookup("json")
	if !ok || jTag == ""{
		return jTag, ok
	}

	//Json tag may contain additional args, return first
	return strings.Split(jTag, " ")[0], true

}