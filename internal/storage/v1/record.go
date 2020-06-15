package v1

import (
	"github.com/aerospike/aerospike-client-go"
	accountv1 "github.com/sajeevany/graphSnapper/internal/account/v1"
	"github.com/sirupsen/logrus"
)

//Record - Aerospike configuration + credentials data
type RecordV1 struct {
	Metadata    Metadata
	Account     Account
	Credentials Credentials
}

func (r RecordV1) ToRecordViewV1() accountv1.RecordViewV1 {
	return accountv1.RecordViewV1{
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
}

func (m Metadata) toMetadataView1() accountv1.MetadataView1 {
	return accountv1.MetadataView1{
		PrimaryKey: m.PrimaryKey,
		LastUpdate: m.LastUpdate,
		CreateTime: m.CreateTime,
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

func (a Account) toAccountView1() accountv1.AccountView1 {
	return accountv1.AccountView1{
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

func (c Credentials) toCredentialsView1() accountv1.CredentialsView1 {
	cv := accountv1.CredentialsView1{
		GrafanaUsers: make(map[string]accountv1.GrafanaUser, len(c.GrafanaUsers)),
	}

	for i, v := range c.GrafanaUsers {
		cv.GrafanaUsers[i] = accountv1.GrafanaUser{
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
		"Metadata",
		map[string]string{
			"PrimaryKey": m.PrimaryKey,
			"LastUpdate": m.LastUpdate,
			"CreateTime": m.CreateTime,
		})
}

func (a Account) getAccountBin() *aerospike.Bin {
	return aerospike.NewBin(
		"Account",
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
		"Credentials",
		map[string]interface{}{
			"GrafanaUsers": grafanaUsersBinMap,
		})
}
