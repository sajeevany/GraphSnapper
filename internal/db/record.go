package db

import (
	"github.com/aerospike/aerospike-client-go"
	"github.com/sirupsen/logrus"
)

type Record interface {
	ToASBinSlice() []*aerospike.Bin
	GetFields() logrus.Fields
}
