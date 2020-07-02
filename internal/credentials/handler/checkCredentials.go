package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sajeevany/graph-snapper/internal/credentials/view"
	"github.com/sajeevany/graph-snapper/internal/db/aerospike/access"
	"github.com/sirupsen/logrus"
	"net/http"
)

const CredGroup = "/credentials"
const TestCredBatch = "/test"

//@Summary CheckCredentials
//@Description Non-authenticated endpoint which tests credentials
//@Produce json
//@Param account body view.AccountViewV1 true "Create account"
//@Success 200 {object} StoredUsers
//@Router /credentials [post]
//@Tags credentials
func CheckCredentials(logger *logrus.Logger, aeroClient *access.ASClient) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		//Bind body to add users object
		var checkCreds view.CheckCredentialsInputView
		if bErr := ctx.ShouldBindJSON(&checkCreds); bErr != nil {
			msg := fmt.Sprintf("Unable to bind request body to CheckCredentialsInputView %v", bErr)
			logger.Errorf(msg)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": msg})
			return
		}

	}
}

func checkCredentials(creds view.CheckCredentialsInputView) {

}