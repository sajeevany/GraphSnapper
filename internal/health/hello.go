package health

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

type Ping struct{
	Response string `json:"response" required:"true" description:"Server hello response" example:"hello"`
}

//@Summary Hello sanity endpoint
//@Description Non-authenticated endpoint that returns 200 with hello message. Used to validate that the service is responsive.
//@Produce json
//@Success 200 {object} health.Ping
//@Router /health/hello [get]
//@Tags health
func Hello(logger *logrus.Logger) gin.HandlerFunc{
	return func(ctx *gin.Context) {
		logger.Println("Hello from within hello")

		//Set response
		ctx.JSON(http.StatusOK, Ping{Response: "hello"})
	}
}
