package credentials

import (
	"fmt"
	"github.com/sirupsen/logrus"
)

type AddUsersModel struct {
	GrafanaUsers []GrafanaUser
}

//IsValid - Returns true if valid. If in valid returns false and an error
func (am AddUsersModel) IsValid() error {

	//Return false if entry is empty
	if len(am.GrafanaUsers) == 0 {
		return fmt.Errorf("GrafanaUsers array is empty")
	}

	//Validate grafana credentials
	for _, user := range am.GrafanaUsers {
		if !user.isValid() {
			return fmt.Errorf("invalid grafana credentials provided key <%v> desc <%v>", user.APIKey, user.Description)
		}
	}

	return nil
}

func (l AddUsersModel) GetFields() logrus.Fields {

	var fields logrus.Fields

	//Get grafana credentials info as fields
	grafanaUserMap := make(map[string]string)
	for _, v := range l.GrafanaUsers {
		grafanaUserMap[v.APIKey] = v.Description
	}
	fields["GrafanaUsers"] = grafanaUserMap

	return fields
}

type GrafanaUser struct {
	APIKey      string
	Description string
}

func (gu GrafanaUser) isValid() bool {
	return gu.APIKey != ""
}

type StoredUsers struct {
	GrafanaUsers []GrafanaDbUser
}

type GrafanaDbUser struct {
	Key         string
	APIKey      string
	Description string
}
