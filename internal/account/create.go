package account

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sajeevany/DockerizedGolangTemplate/internal/db"
	"github.com/sirupsen/logrus"
	"net/http"
)

const AccountGroup = "/account"
const PutAccountEndpoint = "/{id}"

//@Summary Create account record
//@Description Non-authenticated endpoint creates an empty record at the specified key. Overwrites any record that already exists
//@Produce json
//@Param id path string true "id"
//@Success 200 {string} string "ok"
//@Fail 400 {object} gin.H
//@Router /account/{id} [put]
//@Tags account
func AddAccount(logger *logrus.Logger, aeroClient *db.ASClient) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		//Validate that id parameter has been set
		accountId := ctx.Param("id")
		if accountId == "" {
			msg := fmt.Sprintf("Query parameter %v hasn't been set", "id")
			logger.Debug(msg)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": msg})
			return
		}

		ctx.JSON(http.StatusOK, "")
	}
}