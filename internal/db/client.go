package db

import (
	"fmt"
	"github.com/aerospike/aerospike-client-go"
	"github.com/sajeevany/graphSnapper/internal/config"
	"github.com/sirupsen/logrus"
)

type ASClient struct {
	Client           *aerospike.Client
	WritePolicy      *aerospike.WritePolicy
	ScanPolicy       *aerospike.ScanPolicy
	AccountNamespace string
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
		Client:           client,
		WritePolicy:      aerospike.NewWritePolicy(0, 0),
		ScanPolicy:       aerospike.NewScanPolicy(),
		AccountNamespace: conf.AccountNamespace,
	}, nil
}
