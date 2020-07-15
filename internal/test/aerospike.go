package test

import (
	"context"
	"fmt"
	"github.com/docker/go-connections/nat"
	"github.com/sajeevany/graph-snapper/internal/config"
	"github.com/sajeevany/graph-snapper/internal/db/aerospike"
	"github.com/sirupsen/logrus"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"strconv"
	"testing"
)

const (
	aerospikeInternalPort0 = 3000
	aerospikeInternalPort1 = 3001
	aerospikeInternalPort2 = 3002
	aerospikeInternalPort3 = 3003

	aerospikeDockerImageURL = "aerospike:5.0.0.4"
)

func StartAerospikeTestContainer(t *testing.T, ctx context.Context) (testcontainers.Container, *aerospike.ASClient) {

	gRez := testcontainers.ContainerRequest{
		FromDockerfile: testcontainers.FromDockerfile{},
		Image:          aerospikeDockerImageURL,
		ExposedPorts:   getAeroExpostPorts(),
		WaitingFor:     wait.ForLog("service ready: soon there will be cake"),
	}
	aeroContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: gRez,
		Started:          true,
	})
	if err != nil {
		t.Error(err)
	}
	aeroIP, hErr := aeroContainer.Host(ctx)
	if hErr != nil {
		t.Error(hErr)
	}
	ap := strconv.Itoa(aerospikeInternalPort0)
	aeroPort, pErr := aeroContainer.MappedPort(ctx, nat.Port(ap))
	if pErr != nil {
		t.Error(pErr)
	}
	aeroPortInt, sErr := strconv.Atoi(aeroPort.Port())
	if sErr != nil {
		t.Error(sErr)
	}

	//Make aerospike client to be used by test
	asClient, err := aerospike.New(logrus.New(), config.AerospikeCfg{
		Host:                      aeroIP,
		Port:                      aeroPortInt,
		Password:                  "",
		ConnectionRetries:         10,
		ConnectionRetryIntervalMS: 500,
		AccountNamespace: config.AerospikeNamespace{
			Namespace: "test",
			SetName:   "account",
		},
	})
	if err != nil {
		t.Errorf("Failed to create test aerospike client. err <%v>", err)
		t.FailNow()
	}

	return aeroContainer, asClient
}

func getAeroExpostPorts() []string {
	return []string{
		formatTcpPort(aerospikeInternalPort0),
		formatTcpPort(aerospikeInternalPort1),
		formatTcpPort(aerospikeInternalPort2),
		formatTcpPort(aerospikeInternalPort3),
	}
}

func formatTcpPort(port int) string {
	return fmt.Sprintf("%d/tcp", port)
}
