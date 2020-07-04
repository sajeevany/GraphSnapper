package access

import (
	"fmt"
	"github.com/aerospike/aerospike-client-go"
	"github.com/mitchellh/mapstructure"
	"github.com/sajeevany/graph-snapper/internal/db/aerospike/record"
	"strings"
)

type DbReader interface {
	ReadRecord(key *aerospike.Key) (record.Record, error)
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

//KeyExists - Returns true if key exists, with aerospike key and any error that occurs
func (a *AerospikeReader) KeyExists(keyStr string) (bool, *aerospike.Key, error) {

	logger := a.asClient.Logger
	logger.Debugf("Checking if key <%v> exists", keyStr)

	//Create aerospike key to check
	key, err := aerospike.NewKey(a.asClient.AccountNamespace.Namespace, a.asClient.AccountNamespace.SetName, keyStr)
	if err != nil {
		logger.Errorf("Unexpected error when creating new key <%v> ", key.String())
		return false, key, err
	}

	//Check if key exists. Use nil policy because no timeout is required
	exists, kerr := a.asClient.Client.Exists(a.asClient.ReadPolicy, key)
	if kerr != nil {
		logger.Errorf("Error when checking if key <%v>  err <%v>", key.String(), kerr)
		return false, key, kerr
	}
	logger.Debugf("key: <%v> exists: <%v>", key.String(), exists)

	return exists, key, nil
}

func (a *AerospikeReader) ReadRecord(key *aerospike.Key) (record.Record, error) {

	logger := a.asClient.Logger
	aeroClient := a.asClient.Client

	//Get bin map for key
	logger.Debugf("Starting read record for key <%v>", key.String())
	aRecord, rErr := aeroClient.Get(a.asClient.ReadPolicy, key)
	if rErr != nil {
		logger.Errorf("Error when running client.Get operation for key <%v> err <%v>", key.String(), rErr)
	}

	//Get version
	version := record.GetVersion(logger, aRecord.Bins)

	switch strings.ToLower(version) {
	case "":
		vErr := fmt.Errorf("record does not have metadata.version set")
		return nil, vErr
	case record.VersionLevel_1:
		rec, cErr := readV1Record(aRecord.Bins)
		if cErr != nil {
			logger.Errorf("Error converting bin map <%v> to record. err <%v>", aRecord.Bins, cErr)
			return nil, cErr
		}
		logger.WithFields(rec.GetFields()).Debugf("Returning v1 record")
		return rec, nil
	default:
		vErr := fmt.Errorf("record is unsupported version <%v>. update library", version)
		logger.Error(vErr)
		return nil, vErr
	}
}

func readV1Record(bm aerospike.BinMap) (record.Record, error) {

	var rec record.RecordV1
	if cErr := mapstructure.Decode(bm, &rec); cErr != nil {
		return nil, cErr
	}

	return &rec, nil
}
