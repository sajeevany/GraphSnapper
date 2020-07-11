package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sajeevany/graph-snapper/internal/account"
	"github.com/sajeevany/graph-snapper/internal/config"
	"github.com/sajeevany/graph-snapper/internal/credentials"
	"github.com/sajeevany/graph-snapper/internal/db/aerospike"
	"github.com/sajeevany/graph-snapper/internal/health"
	"github.com/sajeevany/graph-snapper/internal/logging"
	"github.com/sajeevany/graph-snapper/internal/logging/middleware"
	"github.com/sirupsen/logrus"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"

	_ "github.com/sajeevany/graph-snapper/docs"
)

const v1Api = "/api/v1"

// @title Graph Snapper API
// @version 1.0
// @description Takes and updates snapshots from a graph service to a document store
// @license.name MIT License
// @BasePath /api/v1
func main() {

	//Create a universal logger. Set default to debug and update later
	logger := logging.Init()
	logger.SetLevel(logrus.DebugLevel)

	//Read configuration file
	confFP := "/app/config/graph-snapper-conf.json"
	conf, isValid, invalidArgs := readConf(logger, confFP)
	if !isValid {
		if prettyIA, err := json.MarshalIndent(invalidArgs, "", "\t"); err != nil {
			logger.WithFields(conf.GetFields()).Fatalf("Configuration file <%v> is invalid. Unable to prettyPrint args <%v>. Invalid arguments: <%v>", confFP, err, invalidArgs)
		} else {
			logger.WithFields(conf.GetFields()).Fatalf("Configuration file <%v> is invalid. Invalid arguments: <%v>", confFP, string(prettyIA))
		}
	}

	//Get aerospike client
	aeroClient, err := aerospike.New(logger, conf.Aerospike)
	if err != nil {
		logger.WithFields(conf.Aerospike.GetFields()).Fatalf("Failed to create Aerospike client using client. Error : <%v>", err)
	}

	//Initialize router
	router := setupRouter(logger)

	//Setup routes
	setupV1Routes(router, logger, aeroClient)

	//Add swagger route
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	//Use default route of 8080.
	routerErr := router.Run(":8080")
	if routerErr != nil {
		logger.Errorf("An error occurred when starting the router. <%v>", routerErr)
	}

}

func readConf(logger *logrus.Logger, filepath string) (*config.Conf, bool, map[string]string) {
	//Read configuration file. Kill startup if an error was found.
	conf, err := config.Read(filepath, logger)
	if err != nil {
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
	engine.Use(middleware.SetCtxLogger(logger))
	engine.Use(middleware.LogRequest(logger))
	engine.Use(gin.Recovery())

	return engine
}

func setupV1Routes(rtr *gin.Engine, logger *logrus.Logger, aeroClient *aerospike.ASClient) {
	addHealthEndpoints(rtr, logger)
	addAccountEndpoints(rtr, logger, aeroClient)
	addCredentialsEndpoints(rtr, logger, aeroClient)
}

func addHealthEndpoints(rtr *gin.Engine, logger *logrus.Logger) {
	v1Api := rtr.Group(fmt.Sprintf("%s%s", v1Api, health.HealthGroup))
	{
		v1Api.GET(health.HelloEndpoint, health.Hello(logger))
	}
}

func addAccountEndpoints(rtr *gin.Engine, logger *logrus.Logger, aeroClient *aerospike.ASClient) {
	v1Api := rtr.Group(fmt.Sprintf("%s%s", v1Api, account.Group))
	{
		v1Api.PUT(account.PutAccountEndpoint, account.PutAccountV1(logger, aeroClient))
		v1Api.GET(account.GetAccountEndpoint, account.GetAccountV1(logger, aeroClient))
	}
}

func addCredentialsEndpoints(rtr *gin.Engine, logger *logrus.Logger, aeroClient *aerospike.ASClient) {
	v1Api := rtr.Group(fmt.Sprintf("%s%s", v1Api, credentials.Group))
	{
		v1Api.POST(credentials.CheckCredentialsEndpoint, credentials.CheckV1(logger))
		v1Api.POST(credentials.AddCredentialsEndpoint, credentials.PutCredentialsV1(logger, aeroClient))
	}
}
