package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sajeevany/graphSnapper/internal/db/aerospike/access"
	"github.com/sajeevany/graphSnapper/internal/db/aerospike/record"
	"github.com/sirupsen/logrus"
	"net/http"
)

const GetAccountEndpoint = "/{id}"

//@Summary Get account record
//@Description Non-authenticated endpoint fetches account at speciied key
//@Produce json
//@Param id path string true "id"
//@Param account body view.Account true "Create account"
//@Success 200 view.Account
//@Fail 404 {object} gin.H
//@Router /account/{id} [get]
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
		record, err := getAccount(logger, aeroClient, accountId, account)
		if err != nil {
			hrErrMsg := fmt.Sprintf("internal error when writing record to aerospike. %v", err)
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

//assumes valid account
func getAccount(logger *logrus.Logger, aeroClient *access.ASClient, key string) (record.RecordV1, error) {

	logger.Debug("Creating account record")

	return record, nil
}
