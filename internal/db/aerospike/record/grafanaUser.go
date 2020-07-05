package record

import "github.com/sirupsen/logrus"

//DBGrafanaUser - Database entry for a GrafanaUser
type DBGrafanaUser struct {
	APIKey      string
	Host        string
	Port        int
	Description string
}

func (u DBGrafanaUser) GetFields() logrus.Fields {
	return logrus.Fields{
		"APIKey":      u.APIKey,
		"Host":        u.Host,
		"Port":        u.Port,
		"Description": u.Description,
	}
}
