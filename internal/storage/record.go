package storage

import (
	accountv1 "github.com/sajeevany/DockerizedGolangTemplate/internal/account/v1"
	"github.com/sirupsen/logrus"
)

type Record interface {
	ToRecordViewV1() accountv1.RecordViewV1
	GetFields() logrus.Fields
}
