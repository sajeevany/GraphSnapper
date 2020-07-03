package credentials

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

const Group = "/credentials"
const CheckCredentialsEndpoint = "check"

//@Summary Check credentials for validaty
//@Description Non-authenticated endpoint Check credentials for validity. Returns an array of user objects with check result
//@Produce json
//@Param credentials body Credentials true "Check credentials"
//@Success 200 {object} CredentialsCheck
//@Fail 400 {object} gin.H
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
		result, err := validateCredentials(creds)
		if err != nil{
			msg := fmt.Sprintf("Error validating credentials. <%v>", err)
			logger.Errorf(msg)
			ctx.JSON(http.StatusInternalServerError, msg)
			return
		}

		ctx.JSON(http.StatusOK,result)
	}
}

func validateCredentials(creds Credentials) (CredentialsCheck, error){

	return CredentialsCheck{}, nil
}