package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sajeevany/graph-snapper/internal/db/aerospike/access"
	"github.com/sajeevany/graph-snapper/internal/db/aerospike/record"
	"github.com/sajeevany/graph-snapper/internal/db/aerospike/view"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

const Group = "/account"
const PutAccountEndpoint = "/{id}"

//@Summary Create account record
//@Description Non-authenticated endpoint creates an empty record at the specified key. Overwrites any record that already exists
//@Produce json
//@Param id path string true "id"
//@Param account body view.AccountViewV1 true "Create account"
//@Success 200 {string} string "ok"
//@Fail 400 {object} gin.H
//@Router /account/{id} [put]
//@Tags account
func PutAccountV1(logger *logrus.Logger, aeroClient *access.ASClient) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		//Validate that id parameter has been set
		accountId := ctx.Param("id")
		if accountId == "" {
			msg := fmt.Sprintf("Query parameter %v hasn't been set", "id")
			logger.Debug(msg)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": msg})
			return
		}

		//Bind account object
		var account view.AccountViewV1
		if bErr := ctx.BindJSON(&account); bErr != nil {
			msg := fmt.Sprintf("Unable to bind request body to account object %v", bErr)
			logger.Errorf(msg)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": msg})
			return
		}

		//Validate account
		if _, vErr := account.IsValid(); vErr != nil {
			logger.WithFields(account.GetFields()).Error("Input account is invalid")
			ctx.JSON(http.StatusBadRequest, gin.H{"error": vErr})
			return
		}

		//Create account
		record, err := createAccount(logger, aeroClient, accountId, account)
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
func createAccount(logger *logrus.Logger, aeroClient *access.ASClient, key string, account view.AccountViewV1) (*record.RecordV1, error) {

	logger.Debug("Creating account record")

	//Get record access known to aerospikeWriter
	rec := newAccountRecordV1(key, account)

	aeroWriter := aeroClient.GetWriter()
	if wErr := aeroWriter.WriteRecord(key, rec); wErr != nil {
		hErr := fmt.Sprintf("Unable to write record with key <%v>", key)
		logger.WithFields(rec.GetFields()).Error(hErr)
		return nil, wErr
	}

	return &rec, nil
}

func newAccountRecordV1(key string, account view.AccountViewV1) record.RecordV1 {
	now := time.Now().UTC().String()
	return record.RecordV1{
		Metadata: record.Metadata{
			PrimaryKey: key,
			LastUpdate: now,
			CreateTime: now,
			Version:    record.V1RecordLevel,
		},
		Account: record.Account{
			Email: account.Email,
			Alias: account.Alias,
		},
		Credentials: record.Credentials{},
	}
}
