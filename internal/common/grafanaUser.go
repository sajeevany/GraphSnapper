package common

import (
	"github.com/sirupsen/logrus"
)

type GrafanaUserV1 struct {
	Auth        Auth
	Description string
}

func (ag GrafanaUserV1) GetFields() logrus.Fields {
	return logrus.Fields{
		"Auth":        ag.Auth.GetRedactedLog(),
		"Description": ag.Description,
	}
}

func (ag GrafanaUserV1) IsValid() bool {
	return ag.Auth.IsValid()
}
