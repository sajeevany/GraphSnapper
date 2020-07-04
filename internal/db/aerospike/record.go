package aerospike

import (
	"github.com/aerospike/aerospike-client-go"
	"github.com/sirupsen/logrus"
)

const (
	MetadataBinName    = "Metadata"
	AccountBinName     = "Account"
	CredentialsBinName = "Credentials"
	VersionAttrName    = "Version"
)

const (
	VersionLevel_1 = "1"
)

type Record interface {
	GetFields() logrus.Fields
	ToASBinSlice() []*aerospike.Bin
	ToRecordViewV1() RecordViewV1
}

//Record - Aerospike configuration + credentials data
type RecordV1 struct {
	Metadata    MetadataV1    `json:"Metadata"`
	Account     AccountV1     `json:"Account"`
	Credentials CredentialsV1 `json:"Credentials"`
}

func (r RecordV1) ToRecordViewV1() RecordViewV1 {
	return RecordViewV1{
		Metadata:    r.Metadata.toMetadataView1(),
		Account:     r.Account.toAccountView1(),
		Credentials: r.Credentials.toCredentialsView1(),
	}
}

//MetadataV1 - Record metadata
type MetadataV1 struct {
	PrimaryKey string
	LastUpdate string
	CreateTime string
	Version    string
}

func (m MetadataV1) toMetadataView1() MetadataViewV1 {
	return MetadataViewV1{
		PrimaryKey:    m.PrimaryKey,
		LastUpdate:    m.LastUpdate,
		CreateTimeUTC: m.CreateTime,
		Version:       m.Version,
	}
}

func (m MetadataV1) GetFields() logrus.Fields {
	return logrus.Fields{
		"PrimaryKey": m.PrimaryKey,
		"LastUpdate": m.LastUpdate,
		"CreateTime": m.CreateTime,
	}
}

//Owner - Creation account details for grouping/fetch
type AccountV1 struct {
	Email string
	Alias string
}

func (a AccountV1) toAccountView1() AccountViewV1 {
	return AccountViewV1{
		Email: a.Email,
		Alias: a.Alias,
	}
}

func (a AccountV1) GetFields() logrus.Fields {
	return logrus.Fields{
		"Email": a.Email,
		"Alias": a.Alias,
	}
}

//CredentialsV1 - CredentialsV1 for various graph and storage services
type CredentialsV1 struct {
	GrafanaUsers map[string]DBGrafanaUser
}

func (c CredentialsV1) toCredentialsView1() CredentialsView1 {
	cv := CredentialsView1{
		GrafanaUsers: make(map[string]GrafanaUser, len(c.GrafanaUsers)),
	}

	for i, v := range c.GrafanaUsers {
		cv.GrafanaUsers[i] = GrafanaUser{
			Description: v.Description,
		}
	}

	return cv
}

func (c CredentialsV1) GetFields() logrus.Fields {
	//Add Grafana user creds
	gFields := logrus.Fields{}
	for i, v := range c.GrafanaUsers {
		gFields[i] = v.GetFields()
	}

	return logrus.Fields{
		"DBGrafanaUser": gFields,
	}
}

//DBGrafanaUser - Database entry for a GrafanaUser
type DBGrafanaUser struct {
	APIKey      string
	Host        string
	Port        int
	Description string
}

func (u DBGrafanaUser) GetFields() logrus.Fields {
	return logrus.Fields{
		"APIKey":      u.APIKey,
		"Host":        u.Host,
		"Port":        u.Port,
		"Description": u.Description,
	}
}

func (r RecordV1) GetFields() logrus.Fields {
	return logrus.Fields{
		"MetadataV1":    r.Metadata.GetFields(),
		"AccountV1":     r.Account.GetFields(),
		"CredentialsV1": r.Credentials.GetFields(),
	}
}

func (r RecordV1) ToASBinSlice() []*aerospike.Bin {
	return []*aerospike.Bin{
		r.Metadata.getMetadataBin(),
		r.Account.getAccountBin(),
		r.Credentials.getCredentialBin(),
	}
}

func (m MetadataV1) getMetadataBin() *aerospike.Bin {
	return aerospike.NewBin(
		MetadataBinName,
		map[string]string{
			"PrimaryKey": m.PrimaryKey,
			"LastUpdate": m.LastUpdate,
			"CreateTime": m.CreateTime,
			"Version":    m.Version,
		})
}

func (a AccountV1) getAccountBin() *aerospike.Bin {
	return aerospike.NewBin(
		AccountBinName,
		map[string]string{
			"Email": a.Email,
			"Alias": a.Alias,
		})
}

func (c CredentialsV1) getCredentialBin() *aerospike.Bin {

	//Create grafana users bin map
	grafanaUsersBinMap := make(map[string]interface{})
	for i, v := range c.GrafanaUsers {
		grafanaUsersBinMap[i] = map[string]string{
			"APIKey":      v.APIKey,
			"Description": v.Description,
		}
	}

	return aerospike.NewBin(
		CredentialsBinName,
		map[string]interface{}{
			"GrafanaUsers": grafanaUsersBinMap,
		})
}
