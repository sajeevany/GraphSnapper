package credentials

import (
	"github.com/sajeevany/graph-snapper/internal/common"
	"github.com/sirupsen/logrus"
)

//AddedCredentials - Users added to specified account
type SetCredentialsV1 struct {
	GrafanaAPIUsers       map[string]common.GrafanaUserV1          `json:"GrafanaAPIUsers"`
	ConfluenceServerUsers map[string]common.ConfluenceServerUserV1 `json:"ConfluenceServerUsers"`
}

func (a SetCredentialsV1) GetFields() logrus.Fields {

	//Convert grafana read users to loggable format
	gru := make(logrus.Fields, len(a.GrafanaAPIUsers))
	for i, v := range a.GrafanaAPIUsers {
		gru[i] = v.GetFields()
	}

	//Convert confluence-server write users to loggable format
	csu := make(logrus.Fields, len(a.GrafanaAPIUsers))
	for i, v := range a.ConfluenceServerUsers {
		csu[i] = v.GetFields()
	}

	return logrus.Fields{
		"GrafanaAPIUsers":       gru,
		"ConfluenceServerUsers": csu,
	}
}

//HasNoUsers - returns true if both grafana and confluence user arrays are empty
func (req SetCredentialsV1) HasNoUsers() bool {
	return len(req.GrafanaAPIUsers) == 0 && len(req.ConfluenceServerUsers) == 0
}
