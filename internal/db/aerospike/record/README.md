#Record design decisions

#####Overall goal:
APIs/handlers should not care about what a record's version. Ideally, records will be of the latest version, but it's not
guaranteed to be as such and so should be updatable by any versioned handler with its respective inputs and behaviour.
It is up to the record to be able to ingest this requirement and be writable in the latest version format. 
	
####Definition:
type **Record** interface {
	GetFields() logrus.Fields
	ToASBinSlice() []*aerospike.Bin
	ToRecordViewV1() RecordViewV1
	AddUserCredentialsV1([]common.GrafanaUserV1, []common.ConfluenceServerUserV1)
}

**GetFields** - Get logrus fields for debug logging in tests

**ToASBinSlice** - Convert record into an Aerospike-compliant format so that it can be easily stored by the record writer 
(record.writer)

**ToRecordViewV1** - Converts a record to a view used to display record information by v1 handlers

**AddUserCredentialsV1** - Updates a record's credentials to include set of IDs defined by a v1 credentials handler

####Conventions going forward:
- Whenever a record is handled, and is determined to be of a non-latest version, then invoke the database writer to rewrite
it in the latest form. Old records should only exist when that particular record is inactive and is not being handled.
- Records should be scheduled to naturally expire so that at some point, inactive records will be pruned and older deprecated
endpoints can be removed
- When a new record version is created, previous record versions (ie RecordV1) should be convertable via an interface method
to the latest version. This will be invoked and by the record writer so that we are always writing in the latest format.
