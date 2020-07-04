package credentials

import (
	"github.com/sajeevany/graph-snapper/internal/logging"
	"github.com/sirupsen/logrus"
)

//CheckCredentialsV1 - Set of credentials to be validated
type CheckCredentialsV1 struct {
	GrafanaReadUsers      []CheckGrafanaReadUserV1
	ConfluenceServerUsers []CheckConfluenceServerUserV1
}

//GrafanaUser - Grafana user with read access
type CheckGrafanaReadUserV1 struct {
	APIKey string
	Host   string
	Port   int
}

func (u CheckGrafanaReadUserV1) GetFields() logrus.Fields {
	return logrus.Fields{
		"APIKey": logging.RedactNonEmpty(u.APIKey),
		"Host":   u.Host,
		"Port":   u.Port,
	}
}

//CheckConfluenceServerUserV1 - confluence user with write access
type CheckConfluenceServerUserV1 struct {
	Username string
	Password string
	Host     string
	Port     int
}

func (u CheckConfluenceServerUserV1) GetFields() logrus.Fields {

	//Redact user and password fields if they have been set
	return logrus.Fields{
		"Username": logging.RedactNonEmpty(u.Username),
		"Password": logging.RedactNonEmpty(u.Password),
		"Host":     u.Host,
		"Port":     u.Port,
	}
}

//CheckCredentialsResultV1 - Check credentials result
type CheckCredentialsResultV1 struct {
	GrafanaReadUserCheck      []CheckGrafanaReadUserResultV1
	ConfluenceServerUserCheck []CheckConfluenceUserResultV1
}

//CheckGrafanaReadUserResultV1 - Grafana read user check result
type CheckGrafanaReadUserResultV1 struct {
	Result bool
	Cause  string `json:"Cause,omitempty"`
	CheckGrafanaReadUserV1
}

//CheckConfluenceUserResultV1 - Confluence write user check result
type CheckConfluenceUserResultV1 struct {
	Result bool
	Cause  string `json:"Cause,omitempty"`
	CheckConfluenceServerUserV1
}
