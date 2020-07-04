package credentials

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

const Group = "/credentials"
const CheckCredentialsEndpoint = "check"

//@Summary Check credentials for validity
//@Description Non-authenticated endpoint Check credentials for validity. Returns an array of user objects with check result
//@Produce json
//@Param credentials body Credentials true "Check credentials"
//@Success 200 {object} CredentialsCheck
//@Fail 400 {object} gin.H
//@Fail 500 {object} gin.H
//@Router /credentials [post]
//@Tags credentials
func CheckV1(logger *logrus.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		logger.Debug("Received check credentials request")

		//Bind credentials object
		var creds Credentials
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

func validateCredentials(logger *logrus.Logger, creds Credentials) (CredentialsCheck, error) {

	logger.Debug("Started credentials validation")
	result := CredentialsCheck{}

	//Check grafana users
	if len(creds.GrafanaReadUsers) != 0 {
		gc, err := authGrafanaUsers(logger, creds.GrafanaReadUsers)
		if err != nil {
			logger.Errorf("Internal error when authenticating grafana users. <%v>", err)
			return result, err
		}
		result.GrafanaReadUserCheck = gc
	}

	//Check confluence users
	if len(creds.ConfluenceServerUsers) != 0 {
		cc, err := authConfluenceUsers(logger, creds.ConfluenceServerUsers)
		if err != nil {
			logger.Errorf("Internal error when authenticating confluence users. <%v>", err)
			return result, err
		}
		result.ConfluenceServerUserCheck = cc
	}

	return result, nil
}
