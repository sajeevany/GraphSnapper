package schedule

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

const (
	Group         = "/schedule"
	CheckEndpoint = "/check"
)

//@Summary Check and execute schedule
//@Description Non-authenticated endpoint which checks and runs a schedule to validate connectivity and storage behaviour by the end user
//@Produce json
//@Param schedule body CheckScheduleV1 true "Check schedule"
//@Success 200 {object} TestResponse
//@Fail 400 {object} gin.H
//@Fail 500 {object} gin.H
//@Router /schedule/check [post]
//@Tags schedule
func CheckV1(logger *logrus.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		logger.Debug("Starting schedule check (v1)")

		//Bind schedule object
		var schedule CheckScheduleV1
		if bErr := ctx.BindJSON(&schedule); bErr != nil {
			msg := fmt.Sprintf("Unable to bind request body to schedule object %v", bErr)
			logger.Errorf(msg)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": msg})
			return
		}

		code, err := checkScheduleV1(schedule)
		if err != nil {
			logger.Error("Received error when running checkScheduleV1. code <%v> err <%v>", code, err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		ctx.Status(code)
	}
}

func checkScheduleV1(schedule CheckScheduleV1) (int, error) {

	//check if input schedule has all required args

	//check if grafana dashboard exists

	//for each panel create snapshot(has expiry)

	//run chromedb login

	//build url to snapshot

	//screen shot each snapshot and save to local dir

	//create page under parent page with correct file names

	//upload all attachments with name to page

	return http.StatusOK, nil
}
