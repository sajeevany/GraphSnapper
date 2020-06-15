package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sajeevany/graphSnapper/internal/db"
	v1 "github.com/sajeevany/graphSnapper/internal/storage/v1"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

const AccountGroup = "/account"
const PutAccountEndpoint = "/{id}"

//@Summary Create account record
//@Description Non-authenticated endpoint creates an empty record at the specified key. Overwrites any record that already exists
//@Produce json
//@Param id path string true "id"
//@Param account body v1.Account true "Create account"
//@Success 200 {string} string "ok"
//@Fail 400 {object} gin.H
//@Router /account/{id} [put]
//@Tags account
func PutAccountV1(logger *logrus.Logger, aeroClient *db.ASClient) gin.HandlerFunc {
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
		var account AccountView1
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
func createAccount(logger *logrus.Logger, aeroClient *db.ASClient, key string, account AccountView1) (*v1.RecordV1, error) {

	logger.Debug("Creating account record")

	//Get record as known to aerospikeWriter
	record := newAccountRecord(key, account)




	return nil, nil
}

func newAccountRecord(key string, account AccountView1) v1.RecordV1{
	now := time.Now().UTC().String()
	return v1.RecordV1{
		Metadata:    v1.Metadata{
			PrimaryKey: key,
			LastUpdate: now,
			CreateTime: now,
		},
		Account:     v1.Account{
			Email: account.Email,
			Alias: account.Alias,
		},
		Credentials: v1.Credentials{},
	}
}
