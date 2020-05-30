package health

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sajeevany/DockerizedGolangTemplate/internal/logging"
	"github.com/sirupsen/logrus"
	"net/http"
)

type Ping struct{
	Response string `json:"response" required:"true" description:"Server hello response" example:"hello"`
}

func Hello(ctx *gin.Context) {
	//Get logger from context
	log := ctx.MustGet(logging.LoggerKey)
	if _, ok := (log).(*logrus.Logger); !ok {
		ctx.JSON(http.StatusInternalServerError, fmt.Errorf("no logger was found using key %v", logging.LoggerKey))
		return
	}

	//Set response
	ctx.JSON(http.StatusOK, Ping{Response: "hello"})
}
