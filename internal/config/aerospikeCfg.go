package config

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"strconv"
)

type AerospikeCfg struct {
	Host             string             `json:"host"`
	Port             int                `json:"port"`
	Password         string             `json:"password"`
	AccountNamespace AerospikeNamespace `json:"accountNamespace"`
}

func (as AerospikeCfg) GetFields() logrus.Fields {
	return logrus.Fields{
		"host":             as.Host,
		"port":             as.Port,
		"password":         as.Password,
		"accountNamespace": as.AccountNamespace.GetFields(),
	}
}

//IsValid - Returns true/false and a non-empty map of all invalid args. Nested args are set in the form of Parent.Child.SubChild
//Inputs:
//    currentPath - json path defined up and including this attribute. ie conf.Aero
//    invalidArgs - map of invalid arguments mapped to their invalid reasons
func (as AerospikeCfg) IsValid(logger *logrus.Logger, currentPath string, invalidArgs map[string]string) bool {

	isValid := true

	//Check attributes
	if as.Host == "" {
		AddInvalidArg(currentPath, "Host", as.Host, invalidArgs)
		isValid = false
	}

	if as.Port <= 0 || as.Port > 65535 {
		AddInvalidArg(currentPath, "Port", strconv.Itoa(as.Port), invalidArgs)
		isValid = false
	}

	//Validate namespace requirements
	if !isAccountNSValid(logger, as.AccountNamespace, currentPath, invalidArgs) {
		isValid = false
	}

	return isValid
}

func isAccountNSValid(logger *logrus.Logger, as AerospikeNamespace, currentPath string, invalidArgs map[string]string) bool {
	isValid := true

	//Check for uninitialzied account namespace
	if (as == AerospikeNamespace{}) {
		AddInvalidArg(currentPath, "AccountNamespace", "", invalidArgs)
		return false
	}

	//Check for invalid account namespace details
	if !as.isValid(logger, fmt.Sprintf("%s.%s", currentPath, "AccountNamespace"), invalidArgs) {
		isValid = false
	}
	return isValid
}

type AerospikeNamespace struct {
	Namespace string `json:"namespace"`
	SetName   string `json:"setName"`
}

func (as AerospikeNamespace) GetFields() logrus.Fields {
	return logrus.Fields{
		"namespace": as.Namespace,
		"setName":   as.SetName,
	}
}

func (as AerospikeNamespace) isValid(logger *logrus.Logger, currentPath string, invalidArgs map[string]string) bool {

	isValid := true

	//Check attributes
	if as.Namespace == "" {
		AddInvalidArg(currentPath, "namespace", as.Namespace, invalidArgs)
		isValid = false
	}

	//Setname is optional. Skip validation

	return isValid
}
