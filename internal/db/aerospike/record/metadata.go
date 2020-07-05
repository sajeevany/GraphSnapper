package record

import (
	"github.com/aerospike/aerospike-client-go"
	"github.com/sirupsen/logrus"
)

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
