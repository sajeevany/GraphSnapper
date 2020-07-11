package aerospike

import (
	"fmt"
	"github.com/aerospike/aerospike-client-go"
	"github.com/sajeevany/graph-snapper/internal/db/aerospike/record"
)

type DbWriter interface {
	WriteRecord(key string, record record.Record) error
	WriteRecordWithASKey(key *aerospike.Key, record record.Record) error
}

func newAerospikeWriter(asClient *ASClient) DbWriter {
	return &AerospikeWriter{
		asClient: asClient,
	}
}

type AerospikeWriter struct {
	asClient *ASClient
}

//Writes record with specified key in the account namespace under the account set. Returns error if one is found
func (a *AerospikeWriter) WriteRecord(key string, record record.Record) error {

	logger := a.asClient.Logger
	logger.WithFields(record.GetFields()).Debug("Starting record create")

	//Create key
	asKey, err := aerospike.NewKey(a.asClient.AccountNamespace.Namespace, a.asClient.AccountNamespace.SetName, key)
	if err != nil {
		logger.Errorf("Unexpected error when creating new key <%v>. err <%v>", key, err)
		return err
	}

	return a.WriteRecordWithASKey(asKey, record)
}

func (a *AerospikeWriter)  WriteRecordWithASKey(asKey *aerospike.Key, record record.Record) error{

	logger := a.asClient.Logger
	logger.WithFields(record.GetFields()).Debug("Starting record create with aerospike key")

	//GetBins
	recBM := record.ToASBinSlice()
	if pErr := a.asClient.Client.PutBins(nil, asKey, recBM...); pErr != nil {
		hErr := fmt.Sprintf("Unable to write bin map <%v> to aerospike namespace <%v> set <%v> key <%v>. err <%v>", recBM, asKey.Namespace(), asKey.SetName(), asKey.String(), pErr)
		logger.WithFields(record.GetFields()).Error(hErr)
		return fmt.Errorf(hErr)
	}

	return nil
}
