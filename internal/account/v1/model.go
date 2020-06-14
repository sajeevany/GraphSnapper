package v1

import (
	"fmt"
	"github.com/sirupsen/logrus"
)

//RecordViewV1 - Aerospike configuration + credentials data
type RecordViewV1 struct {
	Metadata    MetadataView1
	Account     AccountView1
	Credentials CredentialsView1
}

//Metadata - Record metadata
type MetadataView1 struct {
	PrimaryKey string
	LastUpdate string
	CreateTime string
}

//Account - Creation account details
type AccountView1 struct {
	Email string `json:"Email"`
	Alias string `json:"Alias,omitempty"` //Optional arg. Won't be returned if missing.
}

//Credentials - Credentials for various graph and storage services
type CredentialsView1 struct {
	GrafanaUsers map[string]GrafanaUser
}

//GrafanaUser - Grafana user without API key information
type GrafanaUser struct {
	Description string
}

//IsValid - returns true i model is valid. Returns false if invalid and includes a non-nil error
func (a AccountView1) IsValid() (bool, error) {

	if a.Email == "" {
		return false, fmt.Errorf("input email %v is invalid. Expect non-empty value", a.Email)
	}

	//Alias is optional and will not be validated

	return true, nil
}

func (a AccountView1) GetFields() logrus.Fields {
	return logrus.Fields{
		"Email": a.Email,
		"Alias": a.Alias,
	}
}
