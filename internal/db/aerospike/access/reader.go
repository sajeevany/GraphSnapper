package access

import (
	"github.com/aerospike/aerospike-client-go"
	"github.com/sajeevany/graphSnapper/internal/db/aerospike/record"
)

type DbReader interface {
	ReadRecord(key *aerospike.Key) (*record.Record, error)
	KeyExists(key string) (bool, *aerospike.Key, error)
}

func newAerospikeReader(asClient *ASClient) DbReader {
	return &AerospikeReader{
		asClient: asClient,
	}
}

type AerospikeReader struct {
	asClient *ASClient
}

func (a *AerospikeReader) ReadRecord(key *aerospike.Key) (*record.Record, error) {
	return nil, nil
}

//KeyExists - Returns true if key exists, with aerospike key and any error that occurs
func (a *AerospikeReader) KeyExists(keyStr string) (bool, *aerospike.Key, error){

	logger := a.asClient.Logger
	logger.Debugf("Checking if key <%v> exists", keyStr)

	//Create aerospike key to check
	key, err := aerospike.NewKey(a.asClient.AccountNamespace.Namespace, a.asClient.AccountNamespace.SetName, keyStr)
	if err != nil {
		logger.Errorf("Unexpected error when creating new key <%v> namespace <%v> set <%v>", key.String(), key.Namespace(), key.SetName())
		return false, key, err
	}

	//Check if key exists. Use nil policy because no timeout is required
	exists, kerr := a.asClient.Client.Exists(a.asClient.ReadPolicy, key)
	if kerr != nil {
		logger.Errorf("Error when checking if key <%v> namespace <%v> set <%v> exists. err <%v>", key.String(), key.Namespace(), key.SetName(), kerr)
		return false, key, kerr
	}
	logger.Debugf("key: %v exists:%v", key, exists)

	return true, nil, nil
}
