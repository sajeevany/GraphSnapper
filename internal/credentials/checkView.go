package credentials

import (
	"github.com/sajeevany/graph-snapper/internal/logging"
	"github.com/sirupsen/logrus"
)

//Credentials - Set of credentials to be validated
type Credentials struct {
	GrafanaReadUsers      []GrafanaReadUser
	ConfluenceServerUsers []ConfluenceServerUser
}

//GrafanaUser - Grafana user with read access
type GrafanaReadUser struct {
	APIKey string
	Host   string
	Port   int
}

func (u GrafanaReadUser) GetFields() logrus.Fields {
	return logrus.Fields{
		"APIKey": logging.RedactNonEmpty(u.APIKey),
		"Host":   u.Host,
		"Port":   u.Port,
	}
}

//GrafanaUser - Grafana user with read access
type ConfluenceServerUser struct {
	Username string
	Password string
	Host     string
	Port     int
}

func (u ConfluenceServerUser) GetFields() logrus.Fields {

	//Redact user and password fields if they have been set
	return logrus.Fields{
		"Username": logging.RedactNonEmpty(u.Username),
		"Password": logging.RedactNonEmpty(u.Password),
		"Host": u.Host,
		"Port": u.Port,
	}
}



//CredentialsCheck - Check credentials result
type CredentialsCheck struct {
	GrafanaReadUserCheck      []GrafanaReadUserCheck
	ConfluenceServerUserCheck []ConfluenceUserCheck
}

//GrafanaReadUserCheck - Grafana read user check result
type GrafanaReadUserCheck struct {
	Result bool
	Cause  string `json:"Cause,omitempty"`
	GrafanaReadUser
}

//ConfluenceUserCheck - Confluence write user check result
type ConfluenceUserCheck struct {
	Result bool
	Cause  string `json:"Cause,omitempty"`
	ConfluenceServerUser
}
