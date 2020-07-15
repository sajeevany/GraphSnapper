package credentials

import (
	"github.com/sajeevany/graph-snapper/internal/common"
	"github.com/sirupsen/logrus"
)

//CheckCredentialsV1 - Set of credentials to be validated
type CheckCredentialsV1 struct {
	GrafanaReadUsers      []CheckUserV1 `json:"GrafanaAPIUsers"`
	ConfluenceServerUsers []CheckUserV1 `json:"ConfluenceServerUsers"`
}

//GrafanaUser - Grafana user with read access
type CheckUserV1 struct {
	Auth common.Auth
	Host string
	Port int
}

func (u CheckUserV1) GetFields() logrus.Fields {
	return logrus.Fields{
		"Auth": u.Auth.GetFields(),
		"Host": u.Host,
		"Port": u.Port,
	}
}

//CheckCredentialsResultV1 - Check credentials result
type CheckUsersResultV1 struct {
	GrafanaReadUserCheck      []CheckUserResultV1
	ConfluenceServerUserCheck []CheckUserResultV1
}

//CheckUserResultV1 - Grafana read user check result
type CheckUserResultV1 struct {
	Result bool
	Cause  string `json:"Cause,omitempty"`
	CheckUserV1
}
