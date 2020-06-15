package credentials

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sajeevany/graphSnapper/internal/db"
	"github.com/sirupsen/logrus"
	"net/http"
)

const CredGroup = "/credentials"
const PostCredBatch = "/batch"

//@Summary Create users endpoint
//@Description Non-authenticated endpoint which stores the provided credentials credentials owned by a specified account. Verifies connectivity before storing in Aerospike
//@Produce json
//@Param accountId path string true "Account ID"
//@Success 200 {object} StoredUsers
//@Router /credentials/{accountId}/batch [post]
//@Tags credentials
func AddCredentials(logger *logrus.Logger, aeroClient *db.ASClient) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		//Validate that id parameter has been set
		accountId := ctx.Param("accountId")
		if accountId == "" {
			msg := fmt.Sprintf("Query parameter %v hasn't been set", "accountId")
			logger.Debug(msg)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": msg})
			return
		}

		//Bind body to add users object
		var addUsers AddUsersModel
		if bErr := ctx.ShouldBindJSON(&addUsers); bErr != nil {
			msg := fmt.Sprintf("Unable to bind request body to add users model object %v", bErr)
			logger.Errorf(msg)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": msg})
			return
		}

		//Validate request
		if err := addUsers.IsValid(); err != nil {
			logger.Error(err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		//Verify credentials connectivity
		if vErr := validateUserConnectivity(logger, aeroClient, addUsers); vErr != nil {
			logger.WithFields(addUsers.GetFields()).Errorf("Connectivity and credential validation failed")
			ctx.JSON(http.StatusBadRequest, gin.H{"error": vErr.Error()})
			return
		}

		//Store credentials data and return series of keys if no errors are found
		storedUsers, sErr := storeUsers(logger, aeroClient, addUsers)
		if sErr != nil {
			logger.Error("Unable to store users in aerospike", sErr)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": sErr.Error()})
			return
		}

		//Set response
		ctx.JSON(http.StatusOK, storedUsers)
	}
}

func validateUserConnectivity(logger *logrus.Logger, aeroClient *db.ASClient, addUsersMdl AddUsersModel) error {
	return nil
}

func storeUsers(logger *logrus.Logger, aeroClient *db.ASClient, addUsersMdl AddUsersModel) (StoredUsers, error) {

	storedUsers := StoredUsers{
		GrafanaUsers: []GrafanaDbUser{},
	}

	return storedUsers, nil
}
