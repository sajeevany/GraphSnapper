package config

import (
	"github.com/sirupsen/logrus"
	"strconv"
)

type AerospikeCfg struct {
	Host           string `json:"host"`
	Port           int    `json:"port"`
	Password       string `json:"password"`
	GraphNamespace string `json:"graphNamespace"`
}

func (as AerospikeCfg) GetFields() logrus.Fields {
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
//    invalidArgs - map of invalid arguments mapped to their invalid reasons
func (as AerospikeCfg) IsValid(logger *logrus.Logger, currentPath string, invalidArgs map[string]string) bool{

	isValid := true

	//Check attributes
	if as.Host == "" {
		AddInvalidArg(currentPath, "Host", as.Host, invalidArgs)
		isValid = false
	}

	if as.Port <= 0 || as.Port > 65535{
		AddInvalidArg( currentPath, "Port", strconv.Itoa(as.Port), invalidArgs)
		isValid = false
	}

	if as.Password == ""{
		AddInvalidArg(currentPath, "Password", as.Password, invalidArgs)
		isValid = false
	}

	if as.GraphNamespace == ""{
		AddInvalidArg(currentPath, "GraphNamespace", as.GraphNamespace, invalidArgs)
		isValid = false
	}

	return isValid
}