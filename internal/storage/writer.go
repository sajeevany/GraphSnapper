package storage

import (
	"fmt"
	"github.com/aerospike/aerospike-client-go"
	"github.com/sajeevany/graphSnapper/internal/db"
	v1 "github.com/sajeevany/graphSnapper/internal/storage/v1"
	"github.com/sirupsen/logrus"
)

type DbWriter interface {
	WriterRecord(key string, record v1.RecordV1) error
}

func NewAerospikeWriter(logger *logrus.Logger, asClient *db.ASClient) DbWriter {
	return &AerospikeWriter{
		logger:   logger,
		asClient: asClient,
	}
}

type AerospikeWriter struct {
	logger   *logrus.Logger
	asClient *db.ASClient
}

//Writes record with specified key in the account namespace under the account set. Returns error if one is found
func (a *AerospikeWriter) WriterRecord(key string, record v1.RecordV1) error {

	a.logger.WithFields(record.GetFields()).Debug("Starting record create")

	//Create key
	asKey, err := aerospike.NewKey(a.asClient.AccountNamespace, "account", key)
	if err != nil {
		a.logger.Errorf("Unexpected error when creating new key <%v>. err <%v>", key, err)
		return err
	}

	//GetBins
	recBM := record.ToASBinSlice()
	if pErr := a.asClient.Client.PutBins(nil, asKey, recBM...); pErr != nil {
		hErr := fmt.Sprintf("Unable to write bin map <%v> to aerospike namespace <%v> set <%v> key <%v>. err <%v>", recBM, asKey.Namespace(), asKey.SetName(), asKey.String(), pErr)
		a.logger.WithFields(record.GetFields()).Error(hErr)
		return fmt.Errorf(hErr)
	}

	return nil
}
