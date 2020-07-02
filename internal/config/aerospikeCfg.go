package config

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"strconv"
)

type AerospikeCfg struct {
	Host                      string             `json:"host"`
	Port                      int                `json:"port"`
	Password                  string             `json:"password"`
	ConnectionRetries         int                `json:"connectionRetries"`
	ConnectionRetryIntervalMS int                `json:"connectionRetryIntervalMS"`
	AccountNamespace          AerospikeNamespace `json:"accountNamespace"`
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
func (as AerospikeCfg) IsValid(currentPath string, invalidArgs map[string]string) bool {

	isValid := true

	//Check attributes
	if as.Host == "" {
		AddInvalidArgWithCause(currentPath, "Host", as.Host, "value is empty", invalidArgs)
		isValid = false
	}

	if as.Port <= 0 || as.Port > 65535 {
		AddInvalidArgWithCause(currentPath, "Port", strconv.Itoa(as.Port), "value is 0, negative or greater than 65535", invalidArgs)
		isValid = false
	}

	if as.ConnectionRetries <= 0 {
		AddInvalidArgWithCause(currentPath, "ConnectionRetries", strconv.Itoa(as.ConnectionRetries), "value is negative", invalidArgs)
	}

	if as.ConnectionRetryIntervalMS < 0 || as.ConnectionRetryIntervalMS > 10000 {
		AddInvalidArgWithCause(currentPath, "ConnectionRetries", strconv.Itoa(as.ConnectionRetries), "value is negative or exceeds maximum of 10000 milliseconds", invalidArgs)
	}

	//Validate namespace requirements
	if !isAccountNSValid(as.AccountNamespace, currentPath, invalidArgs) {
		isValid = false
	}

	return isValid
}

func isAccountNSValid(as AerospikeNamespace, currentPath string, invalidArgs map[string]string) bool {
	isValid := true

	//Check for uninitialzied account namespace
	if (as == AerospikeNamespace{}) {
		AddInvalidArgWithCause(currentPath, "AccountNamespace", "", "Value is not defined", invalidArgs)
		return false
	}

	//Check for invalid account namespace details
	if !as.IsValid(fmt.Sprintf("%s.%s", currentPath, "AccountNamespace"), invalidArgs) {
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

func (as AerospikeNamespace) IsValid(currentPath string, invalidArgs map[string]string) bool {

	isValid := true

	//Check attributes
	if as.Namespace == "" {
		AddInvalidArgWithCause(currentPath, "Namespace", as.Namespace, "value is empty", invalidArgs)
		isValid = false
	}

	//Setname is optional. Skip validation

	return isValid
}
