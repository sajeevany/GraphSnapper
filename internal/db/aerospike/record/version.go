package record

import (
	"github.com/aerospike/aerospike-client-go"
	"github.com/sirupsen/logrus"
)

const (
	V1RecordLevel = "1"
)

//GetVersion - returns version as a string. Empty if version not found.
func GetVersion(logger *logrus.Logger, aeroRecord aerospike.BinMap) string {

	if aeroRecord == nil || !hasMetadataBin(logger, aeroRecord) {
		logger.Debugf("Aerospike record is empty or missing metadata bin. Returning empty value access version. Binmap <%v>", aeroRecord)
		return ""
	}

	//Get metadata bin map
	mdBin := aeroRecord[MetadataBinName]

	//Get version value
	switch v := mdBin.(type) {
	case map[string]string:
		logger.Debug("Bin map is [str]str. Returning v[Version]")
		return v[VersionAttrName]
	default:
		logger.Debug("Bin map is not [str]str. Returning empty")
		return ""
	}

}

func hasMetadataBin(logger *logrus.Logger, aeroRecord aerospike.BinMap) bool {

	if aeroRecord == nil {
		logger.Debug("Aerospike record is empty. Returning false for hahasMetadataBins")
		return false
	}

	_, ok := aeroRecord[MetadataBinName]

	return ok
}
