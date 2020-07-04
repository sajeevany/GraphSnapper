package credentials

import (
	"github.com/gin-gonic/gin"
	"github.com/sajeevany/graph-snapper/internal/db/aerospike"
	"github.com/sirupsen/logrus"
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
	return func(context *gin.Context) {

	}
}
