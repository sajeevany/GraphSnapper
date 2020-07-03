package credentials

import (
	"github.com/gin-gonic/gin"
	"github.com/sajeevany/graph-snapper/internal/db/aerospike/access"
	"github.com/sirupsen/logrus"
)

const URL = "check"

//@Summary Check credentials for validaty
//@Description Non-authenticated endpoint Check credentials for validity. Returns an array of user objects with check result
//@Produce json
//@Param account body view.AccountViewV1 true "Create account"
//@Success 200 {string} string "ok"
//@Fail 400 {object} gin.H
//@Router /account/:id [put]
//@Tags account
func CheckV1(logger *logrus.Logger, aeroClient *access.ASClient) gin.HandlerFunc {
	return func(ctx *gin.Context) {

	}
}