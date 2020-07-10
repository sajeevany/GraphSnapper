package credentials

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

const (
	Group                    = "/credentials"
	CheckCredentialsEndpoint = "check"
	AddCredentialsEndpoint   = ""
)

//@Summary Check credentials for validity
//@Description Non-authenticated endpoint Check credentials for validity. Returns an array of user objects with check result
//@Produce json
//@Param credentials body CheckCredentialsV1 true "Check credentials"
//@Success 200 {object} CheckUsersResultV1
//@Fail 400 {object} gin.H
//@Fail 500 {object} gin.H
//@Router /credentials/check [post]
//@Tags credentials
func CheckV1(logger *logrus.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		logger.Debug("Received check credentials request")

		//Bind credentials object
		var creds CheckCredentialsV1
		if bErr := ctx.BindJSON(&creds); bErr != nil {
			msg := fmt.Sprintf("Unable to bind request body to credentials object %v", bErr)
			logger.Errorf(msg)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": msg})
			return
		}

		//Validate credentials
		result, err := validateCredentials(logger, creds)
		if err != nil {
			msg := fmt.Sprintf("Error validating credentials. <%v>", err)
			logger.Errorf(msg)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		ctx.JSON(http.StatusOK, result)
	}
}

func validateCredentials(logger *logrus.Logger, creds CheckCredentialsV1) (CheckUsersResultV1, error) {

	logger.Debug("Started credentials validation")
	result := CheckUsersResultV1{}

	//Check grafana users
	if len(creds.GrafanaReadUsers) != 0 {
		result.GrafanaReadUserCheck = authGrafanaUsers(logger, creds.GrafanaReadUsers)
	}

	//Check confluence users
	if len(creds.ConfluenceServerUsers) != 0 {
		result.ConfluenceServerUserCheck = authConfluenceUsers(logger, creds.ConfluenceServerUsers)
	}

	return result, nil
}
