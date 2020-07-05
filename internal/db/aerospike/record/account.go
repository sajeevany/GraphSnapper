package record

import (
	"github.com/aerospike/aerospike-client-go"
	"github.com/sirupsen/logrus"
)

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

func (a AccountV1) getAccountBin() *aerospike.Bin {
	return aerospike.NewBin(
		AccountBinName,
		map[string]string{
			"Email": a.Email,
			"Alias": a.Alias,
		})
}
