package credentials

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sajeevany/graph-snapper/internal/db/aerospike"
	"github.com/sirupsen/logrus"
	"net/http"
)

const PostCredentialsEndpoint = "/{accountID}"

//@Summary Add credentials to an account
//@Description Non-authenticated endpoint that adds grafana and confluence-server users to an account. Assumes entries are pre-validated
//@Produce json
//@Param account body AddCredentialsV1 true "Add credentials"
//@Success 200 {object} AddedCredentialsV1
//@Fail 404 {object} gin.H
//@Fail 500 {object} gin.H
//@Router /credentials [post]
//@Tags credentials
func PostCredentialsV1(logger *logrus.Logger, aeroClient *aerospike.ASClient) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		//Bind add credentials object
		var addReq AddCredentialsV1
		if bErr := ctx.BindJSON(&addReq); bErr != nil {
			msg := fmt.Sprintf("Unable to bind request body to AddCredentialsV1 object %v", bErr)
			logger.Errorf(msg)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": msg})
			return
		}

		//Validate account
		if vErr, returnCode := validateRequest(logger, aeroClient, addReq); vErr != nil {
			logger.WithFields(addReq.GetFields()).Errorf("Input credentials are invalid <%v>", vErr)
			ctx.JSON(returnCode, gin.H{
				"humanReadableError": "Input credentials variables are invalid. Host, user, password, apikey must be non empty. Port must be within ",
				"error":              vErr.Error(),
			})
			return
		}

		//Check if any users have been provided. If not, skip and return 200
		if addReq.HasNoUsers() {
			ctx.Status(http.StatusOK)
			return
		}

		aErr := addUsersToAccount(logger, aeroClient, addReq)
		if aErr != nil {
			hMsg := "Internal error when adding users to data store"
			logger.WithFields(addReq.GetFields()).Error(hMsg, aErr)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"humanReadableError": hMsg,
				"error":              aErr.Error()})
			return
		}
	}
}

//Checks if input is in acceptable and a record exists with the specified key. Returns a non-zero return code if an error is present. Returns no error and a statusOk(200).
func validateRequest(logger *logrus.Logger, aeroClient *aerospike.ASClient, addReq AddCredentialsV1) (error, int) {

	//Validate the account info. Checks if record exists with the ID
	if returnCode, aErr := validateAcctID(logger, aeroClient, addReq.AccountID); aErr != nil {
		return aErr, returnCode
	}

	//Validate grafana users
	for _, gUser := range addReq.GrafanaReadUsers {
		if !gUser.IsValid() {
			logger.WithFields(gUser.GetFields()).Errorf("Grafana user has invalid attributes")
			return fmt.Errorf("grafana user <%#v> is invalid", gUser), http.StatusBadRequest
		}
	}

	//Validate confluence users
	for _, csUser := range addReq.ConfluenceServerUsers {
		if !csUser.IsValid() {
			logger.WithFields(csUser.GetFields()).Errorf("Confluence user has invalid attributes")
			return fmt.Errorf("confluence user <%#v> is invalid", csUser), http.StatusBadRequest
		}
	}

	return nil, http.StatusOK
}

//Returns error if invalid. int value is the http return code to use
func validateAcctID(logger *logrus.Logger, aeroClient *aerospike.ASClient, id string) (int, error) {

	if id == "" {
		return http.StatusBadRequest, fmt.Errorf("account ID is empty and must be defined")
	}

	//check if account exists
	acctExists, _, rErr := aeroClient.GetReader().KeyExists(id)
	if rErr != nil {
		logger.Errorf("Error when reading from aerospike to check if key exists <%v>", rErr)
		return http.StatusInternalServerError, rErr
	}
	if !acctExists {
		return http.StatusNotFound, nil
	}

	return http.StatusOK, nil
}

func addUsersToAccount(logger *logrus.Logger, client *aerospike.ASClient, req AddCredentialsV1) error {
	return nil
}
