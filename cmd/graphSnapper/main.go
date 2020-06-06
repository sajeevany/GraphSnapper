package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sajeevany/DockerizedGolangTemplate/internal/config"
	"github.com/sajeevany/DockerizedGolangTemplate/internal/health"
	"github.com/sajeevany/DockerizedGolangTemplate/internal/logging"
	lm "github.com/sajeevany/DockerizedGolangTemplate/internal/logging/middleware"
	"github.com/sirupsen/logrus"
)

const v1Api = "/api/v1"

func main() {

	//Create a universal logger. Set default to debug and update later
	logger := logging.Init()
	logger.SetLevel(logrus.DebugLevel)

	//Read configuration file
	confFP := "/app/config/graphSnapper-conf.json"
	conf, isValid, invalidArgs := readConf(logger, confFP)
	if !isValid{
		if prettyIA, err := json.MarshalIndent(invalidArgs,"", "\t"); err != nil{
			logger.WithFields(conf.GetFields()).Fatalf("Configuration file <%v> is invalid. Unable to prettyPrint args <%v>. Invalid arguments: <%v>", confFP, err, invalidArgs)
		}else{
			logger.WithFields(conf.GetFields()).Fatalf("Configuration file <%v> is invalid. Invalid arguments: <%v>", confFP, prettyIA)
		}
	}

	//Initialize router
	router := setupRouter(logger)

	//Setup routes
	setupV1Routes(router, logger)

	//Use default route of 8080.
	routerErr := router.Run("8080")
	if routerErr != nil {
		logger.Errorf("An error occurred when starting the router. <%v>", routerErr)
	}

}

func readConf(logger *logrus.Logger, filepath string) (*config.Conf, bool, map[string]string){
	//Read configuration file. Kill startup if an error was found.
	conf, err := config.Read(filepath, logger)
	if err != nil {
		//Log error and use default values returned
		logger.Fatal(err)
	}

	//validate config.
	isValid, invalidArgs := conf.IsValid(logger)

	return conf, isValid, invalidArgs
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

func setupV1Routes(rtr *gin.Engine, logger *logrus.Logger) {
	addHealthEndpoints(rtr, logger)
}

func addHealthEndpoints(rtr *gin.Engine, logger *logrus.Logger) {
	v1 := rtr.Group(fmt.Sprintf("%s%s", v1Api, health.HealthGroup))
	{
		v1.GET(health.HelloEndpoint, health.Hello(logger))
	}
}
