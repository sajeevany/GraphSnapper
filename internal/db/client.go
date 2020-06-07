package db

import (
	"fmt"
	"github.com/aerospike/aerospike-client-go"
	"github.com/sajeevany/DockerizedGolangTemplate/internal/config"
	"github.com/sirupsen/logrus"
)

type ASClient struct {
	Client            *aerospike.Client
	WritePolicy       *aerospike.WritePolicy
	ScanPolicy        *aerospike.ScanPolicy
	GraphNamespace    string
	DocumentNamespace string
}

//New - Returns ASClinet built from config
func New(logger *logrus.Logger, conf config.AerospikeCfg) (*ASClient, error) {

	client, err := aerospike.NewClient(conf.Host, conf.Port)
	if err != nil {
		msg := fmt.Sprintf("Unexpected error when creating aerospike client, <%v> with config.", err)
		logger.WithFields(conf.GetFields()).Error(msg)
		return nil, err
	}
	logger.WithFields(conf.GetFields()).Info("Successful creation of aerospike client")

	//Create policies and define ASClient
	return &ASClient{
		Client:            client,
		WritePolicy:       aerospike.NewWritePolicy(0, 0),
		ScanPolicy:        aerospike.NewScanPolicy(),
		GraphNamespace:    conf.GraphNamespace,
		DocumentNamespace: conf.DocumentNamespace,
	}, nil
}
