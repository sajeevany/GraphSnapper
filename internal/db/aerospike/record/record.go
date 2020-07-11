package record

import (
	"fmt"
	"github.com/aerospike/aerospike-client-go"
	"github.com/lithammer/shortuuid"
	"github.com/sajeevany/graph-snapper/internal/common"
	"github.com/sirupsen/logrus"
)

const (
	MetadataBinName    = "Metadata"
	AccountBinName     = "Account"
	CredentialsBinName = "Credentials"
	VersionAttrName    = "Version"

	GrafanaAPIUserNamespace            = "GrafanaAPIUser"
	ConfluenceServerBasicUserNamespace = "ConfluenceServerBasicUser"
)

const (
	VersionLevel_1 = "1"
)

type Record interface {
	//GetFields - returns logrus fields for logging
	GetFields() logrus.Fields
	//ToASBinSlice - converts record to bin map. Used to write record to db in the latest record format
	ToASBinSlice() []*aerospike.Bin
	//ToRecordViewV1 - converts to v1 record view
	ToRecordViewV1() RecordViewV1
	//SetUserCredentialsV1 - Adds input credentials to record. Does not overwrite any existing records
	SetUserCredentialsV1(*logrus.Logger, []common.GrafanaUserV1, []common.ConfluenceServerUserV1)
}

//Record - Aerospike configuration + credentials data
type RecordV1 struct {
	Metadata    MetadataV1    `json:"Metadata"`
	Account     AccountV1     `json:"Account"`
	Credentials CredentialsV1 `json:"Credentials"`
}

func (r *RecordV1) ToRecordViewV1() RecordViewV1 {
	return RecordViewV1{
		Metadata:    r.Metadata.toMetadataView1(),
		Account:     r.Account.toAccountView1(),
		Credentials: r.Credentials.toCredentialsView1(),
	}
}

func (r *RecordV1) GetFields() logrus.Fields {
	return logrus.Fields{
		"MetadataV1":    r.Metadata.GetFields(),
		"AccountV1":     r.Account.GetFields(),
		"CredentialsV1": r.Credentials.GetFields(),
	}
}

//ToASBinSlice - converts to aerospike bins. Currently writes record in recordv1 format
func (r *RecordV1) ToASBinSlice() []*aerospike.Bin {
	return []*aerospike.Bin{
		r.Metadata.getMetadataBin(),
		r.Account.getAccountBin(),
		r.Credentials.getCredentialBin(),
	}
}

//Add user details to record. Does not overwrite existing users
func (r *RecordV1) SetUserCredentialsV1(logger *logrus.Logger, grafanaUsers []common.GrafanaUserV1, confluenceUsers []common.ConfluenceServerUserV1) {

	logger.Info("Populating record")
	//Add the grafana users
	for _, gu := range grafanaUsers {
		index := getNextFreeGUIdx(r.Credentials.GrafanaAPIUsers, GrafanaAPIUserNamespace)
		logger.Infof("GU got new index <%v>", index)
		r.Credentials.GrafanaAPIUsers[index] = gu
	}

	//Add the confluence users
	for _, cu := range confluenceUsers {
		index := getNextFreeCSUIdx(r.Credentials.ConfluenceServerAPIUsers, ConfluenceServerBasicUserNamespace)
		logger.Infof("CU got new index <%v>", index)
		r.Credentials.ConfluenceServerAPIUsers[index] = cu
	}
	logger.WithFields(r.GetFields()).Info("Record populated")
}

func getNextFreeGUIdx(users map[string]common.GrafanaUserV1, namespace string) string {

	idx := fmt.Sprintf("%s_%s", namespace, shortuuid.New())
	_, keyInUse := users[idx]

	for keyInUse {
		idx = fmt.Sprintf("%s_%s", namespace, shortuuid.New())
		_, keyInUse = users[idx]
		fmt.Println("Going going")
		fmt.Printf("index %v\n",idx)
	}

	return idx
}

func getNextFreeCSUIdx(users map[string]common.ConfluenceServerUserV1, namespace string) string {

	idx := fmt.Sprintf("%s_%s", namespace, shortuuid.New())
	_, keyInUse := users[idx]

	for keyInUse {
		idx = fmt.Sprintf("%s_%s", namespace, shortuuid.New())
		_, keyInUse = users[idx]
	}

	return idx
}
