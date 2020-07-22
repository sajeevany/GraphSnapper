package common

import "github.com/sirupsen/logrus"

type GrafanaDashBoard struct {
	Host   string
	Port   int
	UID    string
	Panels map[string]Panel //if empty include all panels, if non empty only do these panels
	User   GrafanaUserV1
}

func (b GrafanaDashBoard) GetFields() logrus.Fields {

	panels := make(logrus.Fields, len(b.Panels))
	for key, panel := range b.Panels {
		panels[key] = panel.GetFields()
	}

	return logrus.Fields{
		"Host":   b.Host,
		"Port":   b.Port,
		"UID":    b.UID,
		"User":   b.User.GetFields(),
		"Panels": panels,
	}
}

type Panel struct {
	Filler string
}

//GetFields - Filler method for now. At some point user will be able to specify snapshot size for this graph. This is preferable to a version bump
func (p Panel) GetFields() logrus.Fields {
	return logrus.Fields{}
}
