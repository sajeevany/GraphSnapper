package record

import (
	"fmt"
	"github.com/sajeevany/graph-snapper/internal/common"
	"github.com/sirupsen/logrus"
)

//RecordViewV1 - Aerospike configuration + credentials data
type RecordViewV1 struct {
	Metadata    MetadataViewV1   `json:"Metadata"`
	Account     AccountViewV1    `json:"Account"`
	Credentials CredentialsView1 `json:"Credentials"`
}

//MetadataViewV1 - Record metadata
type MetadataViewV1 struct {
	PrimaryKey    string `json:"PrimaryKey"`
	LastUpdate    string `json:"LastUpdate"`
	CreateTimeUTC string `json:"CreateTimeUTC"`
	Version       string `json:"Version"`
}

//AccountViewV1 - Creation account details
type AccountViewV1 struct {
	Email string `json:"Email"`
	Alias string `json:"Alias,omitempty"` //Optional arg. Won't be returned if missing.
}

//Credentials - Credentials for various graph and storage services
type CredentialsView1 struct {
	GrafanaAPIUsers       map[string]GrafanaAPIUser       `json:"GrafanaAPIUsers"`
	ConfluenceServerUsers map[string]ConfluenceServerUser `json:"ConfluenceServerUser"`
}

//GrafanaAPIUser - Grafana user without API key information
type GrafanaAPIUser struct {
	Auth        common.Auth
	Description string
}

type ConfluenceServerUser struct {
	Auth        common.Auth
	Description string
}

//IsValid - returns true if model is valid. Returns false if invalid and includes a non-nil error
func (a AccountViewV1) IsValid() (bool, error) {

	if a.Email == "" {
		return false, fmt.Errorf("input email %v is invalid. Expect non-empty value", a.Email)
	}

	//Alias is optional and will not be validated

	return true, nil
}

func (a AccountViewV1) GetFields() logrus.Fields {
	return logrus.Fields{
		"Email": a.Email,
		"Alias": a.Alias,
	}
}
