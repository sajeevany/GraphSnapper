package record

import (
	"github.com/aerospike/aerospike-client-go"
	"github.com/sajeevany/graphSnapper/internal/db/aerospike/view"
	"github.com/sirupsen/logrus"
)

const (
	MetadataBinName    = "Metadata"
	AccountBinName     = "Account"
	CredentialsBinName = "Credentials"
)

type Record interface {
	GetFields() logrus.Fields
	ToASBinSlice() []*aerospike.Bin
	ToRecordViewV1() view.RecordViewV1
}

//Record - Aerospike configuration + credentials data
type RecordV1 struct {
	Metadata    Metadata
	Account     Account
	Credentials Credentials
}

func (r RecordV1) ToRecordViewV1() view.RecordViewV1 {
	return view.RecordViewV1{
		Metadata:    r.Metadata.toMetadataView1(),
		Account:     r.Account.toAccountView1(),
		Credentials: r.Credentials.toCredentialsView1(),
	}
}

//Metadata - Record metadata
type Metadata struct {
	PrimaryKey string
	LastUpdate string
	CreateTime string
	Version    string
}

func (m Metadata) toMetadataView1() view.MetadataView1 {
	return view.MetadataView1{
		PrimaryKey:    m.PrimaryKey,
		LastUpdate:    m.LastUpdate,
		CreateTimeUTC: m.CreateTime,
		Version:       m.Version,
	}
}

func (m Metadata) GetFields() logrus.Fields {
	return logrus.Fields{
		"PrimaryKey": m.PrimaryKey,
		"LastUpdate": m.LastUpdate,
		"CreateTime": m.CreateTime,
	}
}

//Owner - Creation account details for grouping/fetch
type Account struct {
	Email string
	Alias string
}

func (a Account) toAccountView1() view.AccountViewV1 {
	return view.AccountViewV1{
		Email: a.Email,
		Alias: a.Alias,
	}
}

func (a Account) GetFields() logrus.Fields {
	return logrus.Fields{
		"Email": a.Email,
		"Alias": a.Alias,
	}
}

//Credentials - Credentials for various graph and storage services
type Credentials struct {
	GrafanaUsers map[string]DBGrafanaUser
}

func (c Credentials) toCredentialsView1() view.CredentialsView1 {
	cv := view.CredentialsView1{
		GrafanaUsers: make(map[string]view.GrafanaUser, len(c.GrafanaUsers)),
	}

	for i, v := range c.GrafanaUsers {
		cv.GrafanaUsers[i] = view.GrafanaUser{
			Description: v.Description,
		}
	}

	return cv
}

func (c Credentials) GetFields() logrus.Fields {
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
	Description string
}

func (u DBGrafanaUser) GetFields() logrus.Fields {
	return logrus.Fields{
		"APIKey":      u.APIKey,
		"Description": u.Description,
	}
}

func (r RecordV1) GetFields() logrus.Fields {
	return logrus.Fields{
		"Metadata":    r.Metadata.GetFields(),
		"Account":     r.Account.GetFields(),
		"Credentials": r.Credentials.GetFields(),
	}
}

func (r RecordV1) ToASBinSlice() []*aerospike.Bin {
	return []*aerospike.Bin{
		r.Metadata.getMetadataBin(),
		r.Account.getAccountBin(),
		r.Credentials.getCredentialBin(),
	}
}

func (m Metadata) getMetadataBin() *aerospike.Bin {
	return aerospike.NewBin(
		MetadataBinName,
		map[string]string{
			"PrimaryKey": m.PrimaryKey,
			"LastUpdate": m.LastUpdate,
			"CreateTime": m.CreateTime,
			"Version":    m.Version,
		})
}

func (a Account) getAccountBin() *aerospike.Bin {
	return aerospike.NewBin(
		AccountBinName,
		map[string]string{
			"Email": a.Email,
			"Alias": a.Alias,
		})
}

func (c Credentials) getCredentialBin() *aerospike.Bin {

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
