package access

import (
	"fmt"
	"github.com/aerospike/aerospike-client-go"
	"github.com/sajeevany/graphSnapper/internal/config"
	"github.com/sirupsen/logrus"
)

type AerospikeClient interface {
	GetWriter() DbWriter
	GetReader() DbReader
}

type ASClient struct {
	Logger           *logrus.Logger
	Client           *aerospike.Client
	WritePolicy      *aerospike.WritePolicy
	ReadPolicy       *aerospike.BasePolicy
	AccountNamespace config.AerospikeNamespace
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
		Logger:           logger,
		Client:           client,
		WritePolicy:      aerospike.NewWritePolicy(0, 0),
		ReadPolicy:       aerospike.NewPolicy(),
		AccountNamespace: conf.AccountNamespace,
	}, nil
}

func (a *ASClient) GetWriter() DbWriter {
	return NewAerospikeWriter(a)
}

func (a *ASClient) GetReader() DbReader {
	return NewAerospikeReader(a)
}
