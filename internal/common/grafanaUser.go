package common

import (
	"github.com/sajeevany/graph-snapper/internal/config"
	"github.com/sirupsen/logrus"
)

type GrafanaUserV1 struct {
	Auth        Auth
	Host        string
	Port        int
	Description string
}

func (ag GrafanaUserV1) GetFields() logrus.Fields {
	return logrus.Fields{
		"Auth":        ag.Auth.GetRedactedLog(),
		"Host":        ag.Host,
		"Port":        ag.Port,
		"Description": ag.Description,
	}
}

func (ag GrafanaUserV1) IsValid() bool {
	return ag.Auth.IsValid() && ag.Host != "" && config.IsPortValid(ag.Port)
}
