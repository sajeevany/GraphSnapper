package record

import (
	"github.com/aerospike/aerospike-client-go"
	"github.com/lithammer/shortuuid"
	"github.com/sajeevany/graph-snapper/internal/common"
	"github.com/sirupsen/logrus"
)

//CredentialsV1 - CredentialsV1 for various graph and storage services
type CredentialsV1 struct {
	GrafanaAPIUsers map[string]common.GrafanaUserV1
}

func (c CredentialsV1) toCredentialsView1() CredentialsView1 {
	cv := CredentialsView1{
		GrafanaUsers: make(map[string]GrafanaUser, len(c.GrafanaAPIUsers)),
	}

	for i, v := range c.GrafanaAPIUsers {
		cv.GrafanaUsers[i] = GrafanaUser{
			Description: v.Description,
		}
	}

	return cv
}

func (c CredentialsV1) GetFields() logrus.Fields {
	//Add Grafana user creds
	gFields := logrus.Fields{}
	for i, v := range c.GrafanaAPIUsers {
		gFields[i] = v.GetFields()
	}

	return logrus.Fields{
		"DBGrafanaUser": gFields,
	}
}

func (c CredentialsV1) getCredentialBin() *aerospike.Bin {

	//Create grafana users bin map
	grafanaUsersBinMap := make(map[string]interface{})
	for i, v := range c.GrafanaAPIUsers {
		grafanaUsersBinMap[i] = map[string]interface{}{
			"Authentication":      v.Authentication.ToAerospikeBinMap(),
			"Description": v.Description,
		}
	}

	return aerospike.NewBin(
		CredentialsBinName,
		map[string]interface{}{
			"GrafanaAPIUsers": grafanaUsersBinMap,
		})
}

func getNextFreeIdx(users map[string]common.GrafanaUserV1, namespace string) string {

	idx := shortuuid.NewWithNamespace(namespace)
	_, keyInUse := users[idx]

	for keyInUse {
		idx = shortuuid.NewWithNamespace(namespace)
		_, keyInUse = users[idx]
	}

	return idx
}
