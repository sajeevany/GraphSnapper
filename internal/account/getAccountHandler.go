package account

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sajeevany/graph-snapper/internal/db/aerospike"
	"github.com/sirupsen/logrus"
	"net/http"
)

const GetAccountEndpoint = "/:id"

//@Summary Get account record
//@Description Non-authenticated endpoint fetches account at specified key
//@Produce json
//@Param id path string true "id"
//@Success 200 {object} aerospike.RecordViewV1
//@Fail 404 {object} gin.H
//@Router /account/:id [get]
//@Tags account
func GetAccountV1(logger *logrus.Logger, aeroClient *aerospike.ASClient) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		//Validate that id parameter has been set
		accountId := ctx.Param("id")
		if accountId == "" {
			msg := fmt.Sprintf("Query parameter %v hasn't been set", "id")
			logger.Debug(msg)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": msg})
			return
		}

		//Check if record exists
		reader := aeroClient.GetReader()
		recordExists, aKey, kErr := reader.KeyExists(accountId)
		if kErr != nil {
			hrErrMsg := fmt.Sprintf("unable to check db for key <%v>. err <%v>", accountId, kErr)
			logger.Errorf(hrErrMsg)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error":              kErr,
				"humanReadableError": hrErrMsg,
			})
			return
		}
		if !recordExists {
			logger.Debugf("key <%v> namespace <%v> set <%v> does not exist. Returning 404", accountId, aKey.Namespace(), aKey.SetName())
			ctx.Status(http.StatusNotFound)
			return
		}

		//fetch account
		rec, rErr := reader.ReadRecord(aKey)
		if rErr != nil {
			hrErrMsg := fmt.Sprintf("unable to read db for key <%v> namespace <%v> set <%v>. err <%v>", accountId, aKey.Namespace(), aKey.SetName(), rErr)
			logger.Errorf(hrErrMsg)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error":              rErr,
				"humanReadableError": hrErrMsg,
			})
			return
		}

		//Return view
		view := rec.ToRecordViewV1()
		logger.Infof("Record <%v> as view <%v>", rec, view)
		ctx.JSON(http.StatusOK, view)
	}
}
