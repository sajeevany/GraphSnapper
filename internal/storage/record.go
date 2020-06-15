package storage

import (
	"github.com/aerospike/aerospike-client-go"
	accountv1 "github.com/sajeevany/graphSnapper/internal/account/v1"
	"github.com/sirupsen/logrus"
)

type Record interface {
	ToASBinSlice() []aerospike.Bin
	ToRecordViewV1() accountv1.RecordViewV1
	GetFields() logrus.Fields
}
