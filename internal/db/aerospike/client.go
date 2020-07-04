package aerospike

import (
	"fmt"
	"github.com/aerospike/aerospike-client-go"
	"github.com/sajeevany/graph-snapper/internal/config"
	"github.com/sirupsen/logrus"
	"time"
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

	logger.Debug("Starting aerospike client creation")
	client, err := getAerospikeClient(logger, conf.Host, conf.Port, conf.ConnectionRetries, conf.ConnectionRetryIntervalMS)
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

func getAerospikeClient(logger *logrus.Logger, host string, port, retryTimes int, retryIntervalMilliseconds int) (*aerospike.Client, error) {

	var client *aerospike.Client
	var err error

	//Will attempt to get client at least once
	retryTimes = 1 + abs(retryTimes)
	retryInterval := time.Duration(retryIntervalMilliseconds) * time.Millisecond

	for i := 0; i < retryTimes; i++ {
		logger.Debugf("Connection creation attempt #<%v>", i)
		client, err = aerospike.NewClient(host, port)
		if err == nil {
			return client, err
		}
		//pause between client fetch times
		time.Sleep(retryInterval)
		retryIntervalMilliseconds = retryIntervalMilliseconds * retryIntervalMilliseconds
	}

	logger.Debug("Out of retry attempts ")
	return client, err
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func (a *ASClient) GetWriter() DbWriter {
	return newAerospikeWriter(a)
}

func (a *ASClient) GetReader() DbReader {
	return newAerospikeReader(a)
}
