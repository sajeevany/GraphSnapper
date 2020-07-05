package credentials

import (
	"github.com/sajeevany/graph-snapper/internal/config"
	"github.com/sajeevany/graph-snapper/internal/logging"
	"github.com/sirupsen/logrus"
)

//AddCredentialsV1 - Request to add users to account with specified ID
type AddCredentialsV1 struct {
	AccountID             string
	GrafanaReadUsers      []AddGrafanaReadUserV1      `json:"GrafanaReadUsers"`
	ConfluenceServerUsers []AddConfluenceServerUserV1 `json:"ConfluenceServerUsers"`
}

func (a AddCredentialsV1) GetFields() logrus.Fields {

	//Convert grafana read users to loggable format
	gru := make([]logrus.Fields, len(a.GrafanaReadUsers))
	for i, v := range a.GrafanaReadUsers {
		gru[i] = v.GetFields()
	}

	//Convert confluence-server write users to loggable format
	csu := make([]logrus.Fields, len(a.GrafanaReadUsers))
	for i, v := range a.ConfluenceServerUsers {
		csu[i] = v.GetFields()
	}

	return logrus.Fields{
		"AccountID":             a.AccountID,
		"GrafanaReadUsers":      gru,
		"ConfluenceServerUsers": csu,
	}
}

//HasNoUsers - returns true if both grafana and confluence user arrays are empty
func (req AddCredentialsV1) HasNoUsers() bool {
	return len(req.GrafanaReadUsers) == 0 && len(req.ConfluenceServerUsers) == 0
}

//AddGrafanaReadUserV1 - Grafana user with read access
type AddGrafanaReadUserV1 struct {
	APIKey      string
	Host        string
	Port        int
	Description string
}

func (ag AddGrafanaReadUserV1) GetFields() logrus.Fields {
	return logrus.Fields{
		"APIKey":      logging.RedactNonEmpty(ag.APIKey),
		"Host":        ag.Host,
		"Port":        ag.Port,
		"Description": ag.Description,
	}
}

func (ag AddGrafanaReadUserV1) IsValid() bool {
	return ag.APIKey != "" && ag.Host != "" && config.IsPortValid(ag.Port)
}

//AddConfluenceServerUserV1 - confluence user with write access
type AddConfluenceServerUserV1 struct {
	Username    string
	Password    string
	Host        string
	Port        int
	Description string
}

func (ag AddConfluenceServerUserV1) GetFields() logrus.Fields {
	return logrus.Fields{
		"Username":    logging.RedactNonEmpty(ag.Username),
		"Password":    logging.RedactNonEmpty(ag.Password),
		"Host":        ag.Host,
		"Port":        ag.Port,
		"Description": ag.Description,
	}
}

func (acs AddConfluenceServerUserV1) IsValid() bool {
	return acs.Username != "" && acs.Password != "" && acs.Host != "" && config.IsPortValid(acs.Port)
}

//AddedCredentials - Users added to specified account
type AddedCredentialsV1 struct {
	GrafanaReadUsers      map[string]AddGrafanaReadUserV1      `json:"GrafanaReadUsers"`
	ConfluenceServerUsers map[string]AddConfluenceServerUserV1 `json:"ConfluenceServerUsers"`
}
