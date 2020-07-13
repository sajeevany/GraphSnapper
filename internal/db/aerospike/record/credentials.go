package record

import (
	"github.com/aerospike/aerospike-client-go"
	"github.com/sajeevany/graph-snapper/internal/common"
	"github.com/sirupsen/logrus"
)

const (
	GrafanaAPIUsersBMKey    = "GrafanaAPIUsers"
	ConfluenceAPIUsersBMKey = "ConfluenceServerAPIUsers"
)

//CredentialsV1 - CredentialsV1 for various graph and storage services
type CredentialsV1 struct {
	GrafanaAPIUsers          map[string]common.GrafanaUserV1
	ConfluenceServerAPIUsers map[string]common.ConfluenceServerUserV1
}

func (c CredentialsV1) toCredentialsView1() CredentialsView1 {
	cv := CredentialsView1{
		GrafanaUsers:          make(map[string]GrafanaUser, len(c.GrafanaAPIUsers)),
		ConfluenceServerUsers: make(map[string]ConfluenceServerUser, len(c.ConfluenceServerAPIUsers)),
	}

	for i, v := range c.GrafanaAPIUsers {
		cv.GrafanaUsers[i] = GrafanaUser{
			Auth:        v.Auth.GetRedactedView(),
			Host:        v.Host,
			Port:        v.Port,
			Description: v.Description,
		}
	}

	for i, v := range c.ConfluenceServerAPIUsers {
		cv.ConfluenceServerUsers[i] = ConfluenceServerUser{
			Auth:        v.Auth.GetRedactedView(),
			Host:        v.Host,
			Port:        v.Port,
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

	//Add Confluence user creds
	csFields := logrus.Fields{}
	for i, v := range c.ConfluenceServerAPIUsers {
		csFields[i] = v.GetFields()
	}

	return logrus.Fields{
		GrafanaAPIUsersBMKey:    gFields,
		ConfluenceAPIUsersBMKey: csFields,
	}
}

func (c CredentialsV1) getCredentialBin() *aerospike.Bin {

	//Create grafana users bin map
	grafanaUsersBinMap := make(map[string]interface{})
	for i, v := range c.GrafanaAPIUsers {
		grafanaUsersBinMap[i] = map[string]interface{}{
			"Auth":        v.Auth.ToAerospikeBinMap(),
			"Host":        v.Host,
			"Port":        v.Port,
			"Description": v.Description,
		}
	}

	//Create confluence server users bin map
	confluenceServerUsersBinMap := make(map[string]interface{})
	for i, v := range c.ConfluenceServerAPIUsers {
		confluenceServerUsersBinMap[i] = map[string]interface{}{
			"Auth":        v.Auth.ToAerospikeBinMap(),
			"Host":        v.Host,
			"Port":        v.Port,
			"Description": v.Description,
		}
	}

	return aerospike.NewBin(
		CredentialsBinName,
		map[string]interface{}{
			GrafanaAPIUsersBMKey:    grafanaUsersBinMap,
			ConfluenceAPIUsersBMKey: confluenceServerUsersBinMap,
		})
}
