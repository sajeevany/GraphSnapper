package common

import "github.com/sirupsen/logrus"

type GrafanaDashBoard struct {
	Host             string
	Port             int
	UID              string
	IncludePanelsIDs []string //blank means include all panels. Will include newly added panels
	ExcludePanelsIDs []string //blank means exclude nothing
	User             GrafanaUserV1
}

func (b GrafanaDashBoard) GetFields() logrus.Fields {

	return logrus.Fields{
		"Host":             b.Host,
		"Port":             b.Port,
		"UID":              b.UID,
		"User":             b.User.GetFields(),
		"IncludePanelsIDs": b.IncludePanelsIDs,
		"ExcludePanelsIDs": b.ExcludePanelsIDs,
	}
}
