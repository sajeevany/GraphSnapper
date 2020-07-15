package common

import (
	"github.com/sajeevany/graph-snapper/internal/config"
	"github.com/sirupsen/logrus"
)

type ConfluenceServerUserV1 struct {
	Host        string
	Port        int
	Description string
	Auth        Auth
}

//ConfluenceServerUserV1 - Confluence server user
func (u ConfluenceServerUserV1) GetFields() logrus.Fields {
	return logrus.Fields{
		"Host":        u.Host,
		"Port":        u.Port,
		"Description": u.Description,
		"Auth":        u.Auth.GetFields(),
	}
}

func (acs ConfluenceServerUserV1) IsValid() bool {
	return acs.Auth.IsValid() && acs.Host != "" && config.IsPortValid(acs.Port)
}
