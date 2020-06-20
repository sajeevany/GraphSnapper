package v1

import (
	"fmt"
	"github.com/aerospike/aerospike-client-go"
	"github.com/gin-gonic/gin"
	"github.com/sajeevany/graphSnapper/internal/db/aerospike/access"
	"github.com/sajeevany/graphSnapper/internal/db/aerospike/record"
	"github.com/sirupsen/logrus"
	"net/http"
)

const GetAccountEndpoint = "/:id"

//@Summary Get account record
//@Description Non-authenticated endpoint fetches account at specified key
//@Produce json
//@Param id path string true "id"
//@Success 200 {object} view.RecordViewV1
//@Fail 404 {object} gin.H
//@Router /account/:id [get]
//@Tags account
func GetAccountV1(logger *logrus.Logger, aeroClient *access.ASClient) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		//Validate that id parameter has been set
		accountId := ctx.Param("id")
		if accountId == "" {
			msg := fmt.Sprintf("Query parameter %v hasn't been set", "id")
			logger.Debug(msg)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": msg})
			return
		}

		//Create account
		record, rExists, err := getAccount(logger, aeroClient, accountId)
		if err != nil {
			hrErrMsg := fmt.Sprintf("internal error when writing record to aerospike. %v", err)
			logger.WithFields(record.GetFields()).Errorf(hrErrMsg)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error":              err,
				"humanReadableError": hrErrMsg,
			})
		}

		if !rExists {
			hrErrMsg := fmt.Sprintf("record for key <%v> doesn't exist. err %v", accountId, err)
			logger.WithFields(record.GetFields()).Errorf(hrErrMsg)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error":              err,
				"humanReadableError": hrErrMsg,
			})
		}

		view := record.ToRecordViewV1()

		ctx.JSON(http.StatusOK, view)
	}
}

func recordExists(logger *logrus.Logger, aeroClient *access.ASClient, key string) (bool, *aerospike.Key, error) {
	logger.Debugf("Checking if key <%v> exists", key)

	////Create aerospike key to check
	//key, err := aerospike.NewKey(aeroClient.AccountNamespace, client.SetMetadata.SetName, id)
	//if err != nil {
	//	logger.Errorf("Unexpected error when creating new key <%v>", key)
	//	return true, key, err
	//}
	//
	////Check if key exists. Use nil policy because no timeout is required
	//exists, kerr := aeroClient.Client.Exists(nil, key)
	//if kerr != nil {
	//	logger.Error("Error when checking if key exists", kerr)
	//	return true, key, kerr
	//}
	//logger.Debugf("key: %v exists:%v", key, exists)
	return false, nil, nil
}

//assumes valid account
func getAccount(logger *logrus.Logger, aeroClient *access.ASClient, key string) (*record.RecordV1, bool, error) {

	logger.Debug("Creating account record")

	return &record.RecordV1{}, false, nil
}
