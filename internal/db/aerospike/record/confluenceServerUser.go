package record

import "github.com/sirupsen/logrus"

type DBConfluenceServerUser struct {
	Host           string
	Port           int
	Description    string
	Authentication DBAuth
}

//DBConfluenceServerBasicUser - Database entry for a Confluence server user
func (u DBConfluenceServerUser) GetFields() logrus.Fields {
	return logrus.Fields{
		"Host":           u.Host,
		"Port":           u.Port,
		"Description":    u.Description,
		"Authentication": u.Authentication.GetFields(),
	}
}
