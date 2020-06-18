package access

import "github.com/sajeevany/graphSnapper/internal/db/aerospike/record"

type DbReader interface {
	ReadRecord(key string) (bool, *record.Record, error)
}

func NewAerospikeReader(asClient *ASClient) DbReader {
	return &AerospikeReader{
		asClient: asClient,
	}
}

type AerospikeReader struct {
	asClient *ASClient
}

func (a *AerospikeReader) ReadRecord(key string) (bool, *record.Record, error) {
	return false, nil, nil
}
