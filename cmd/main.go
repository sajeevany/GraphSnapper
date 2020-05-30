package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sajeevany/DockerizedGolangTemplate/internal/config"
	"github.com/sajeevany/DockerizedGolangTemplate/internal/endpoints"
	"github.com/sajeevany/DockerizedGolangTemplate/internal/logging"
	lm "github.com/sajeevany/DockerizedGolangTemplate/internal/logging/middleware"
	"github.com/sirupsen/logrus"
)

const v1Api = "/api/v1"

func main() {

	//Create a universal logger
	logger := logging.Init()

	//Read configuration file
	conf, err := config.Read("/config/default.json", logger)
	if err != nil{
		//Log error and use default values returned
		logger.Error(err)
	}

	//Initialize router
	router := setupRouter(logger)

	//Setup routes
	setupV1Routes(router)

	//Use default route of 8080.
	port := fmt.Sprintf(":%d", conf.Port)
	routerErr := router.Run(port)
	if routerErr != nil {
		logger.Errorf("An error occurred when starting the router. <%v>", routerErr)
	}

}

//setupRouter - Create the router and set middleware
func setupRouter(logger *logrus.Logger) *gin.Engine {

	engine := gin.New()

	//Add middleware
	engine.Use(lm.SetCtxLogger(logger))
	engine.Use(lm.LogRequest(logger))
	engine.Use(gin.Recovery())

	return engine
}

func setupV1Routes(rtr *gin.Engine) {
	addHealthEndpoints(rtr)
}

func addHealthEndpoints(rtr *gin.Engine) {
	v1 := rtr.Group(fmt.Sprintf("%s%s", v1Api, endpoints.HealthGroup))
	{
		hello := endpoints.BuildHelloEndpoint()
		v1.GET(hello.URL, hello.Handlers...)
	}
}
