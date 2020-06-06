package user

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

const UserGroup = "/user"
const PostBatchUsers = "/batch"

//@Summary Create users endpoint
//@Description Non-authenticated endpoint which stores the provided user. Verifies connectivity before storing into Aerospike
//@Produce json
//@Success 200 {object} StoredUsers
//@Router /user [post]
//@Tags user
func CreateUsers(logger *logrus.Logger) gin.HandlerFunc{
	return func(ctx *gin.Context) {

		//Bind body to add users object
		var addUsers AddUsersModel
		if bErr := ctx.ShouldBindJSON(&addUsers); bErr != nil{
			msg := fmt.Sprintf("Unable to bind request body to add users model object %v", bErr)
			logger.Errorf(msg)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": msg})
			return
		}

		//Validate request
		if err := addUsers.IsValid(); err != nil{
			logger.Error(err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}


		//Verify user connectivity

		//

		//Set response
		ctx.JSON(http.StatusOK, StoredUsers{})
	}
}
