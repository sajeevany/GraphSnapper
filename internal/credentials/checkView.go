package credentials

import "github.com/sirupsen/logrus"

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
		"APIKey": u.APIKey,
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
