package view

import (
	"fmt"
	"github.com/sirupsen/logrus"
)

//RecordViewV1 - Aerospike configuration + credentials data
type RecordViewV1 struct {
	Metadata    MetadataView1    `json:"Metadata"`
	Account     AccountViewV1    `json:"Account"`
	Credentials CredentialsView1 `json:"Credentials"`
}

//MetadataView1 - Record metadata
type MetadataView1 struct {
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
	GrafanaUsers map[string]GrafanaUser `json:"GrafanaUsers"`
}

//GrafanaUser - Grafana user without API key information
type GrafanaUser struct {
	Description string `json:"Description"`
}

//IsValid - returns true i model is valid. Returns false if invalid and includes a non-nil error
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
