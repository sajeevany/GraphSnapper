package config

import "github.com/sirupsen/logrus"

type Conf struct {
	Aerospike AerospikeCfg `json:"aerospike"`
}

func (c Conf) getFields() logrus.Fields {
	return logrus.Fields{
		"aerospike": c.Aerospike.getFields(),
	}
}

//IsValid - Returns true/false and a non-empty map of all invalid args. Nested args are set in the form of Parent.Child.SubChild
func (c Conf) IsValid(logger *logrus.Logger) (bool, map[string]string){

	return c.Aerospike.IsValid(logger, "conf.aerospike")
}