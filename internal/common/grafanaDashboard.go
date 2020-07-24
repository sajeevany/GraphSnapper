package common

import "github.com/sirupsen/logrus"

type GrafanaDashBoard struct {
	Host             string
	Port             int
	UID              string
	IncludePanelsIDs []int //blank means include all panels. Will include newly added panels
	ExcludePanelsIDs []int //blank means exclude nothing. New panels will be automatically included
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
