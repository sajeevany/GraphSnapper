package config

import "github.com/sirupsen/logrus"

type Conf struct {
	Aerospike AerospikeCfg `json:"aerospike"`
	Logging   Logging      `json:"logging"`
}

func NewConfWithDefaults() Conf {
	return Conf{
		Aerospike: AerospikeCfg{
			ConnectionRetries:         3,
			ConnectionRetryIntervalMS: 10,
		},
		Logging: Logging{
			Level: "debug",
		},
	}
}

func (c Conf) GetFields() logrus.Fields {
	return logrus.Fields{
		"aerospike": c.Aerospike.GetFields(),
	}
}

//IsValid - Returns true/false and a non-empty map of all invalid args. Nested args are set in the form of Parent.Child.SubChild
func (c Conf) IsValid(logger *logrus.Logger) (bool, map[string]string) {

	var invalidArgs = make(map[string]string)

	aeroIsValid := c.Aerospike.IsValid("conf.aerospike", invalidArgs)
	logIsValid := c.Logging.IsValid("conf.logging", invalidArgs)

	return aeroIsValid && logIsValid, invalidArgs
}
