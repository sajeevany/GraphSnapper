package common

import (
	"github.com/sirupsen/logrus"
)

type ConfluenceServerUserV1 struct {
	Auth        Auth
	Description string
}

//ConfluenceServerUserV1 - Confluence server user
func (u ConfluenceServerUserV1) GetFields() logrus.Fields {
	return logrus.Fields{
		"Description": u.Description,
		"Auth":        u.Auth.GetFields(),
	}
}

func (acs ConfluenceServerUserV1) IsValid() bool {
	return acs.Auth.IsValid()
}
