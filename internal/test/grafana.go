package test

import (
	"context"
	"fmt"
	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"strconv"
	"testing"
)

const (
	//Predefined container API keys
	GrafanaAdminUserAPIKey = "eyJrIjoiUDI4NFg3WFpoNDJXRnJ4RTBZeDMwb2U0WnA2QTBkRmoiLCJuIjoiYWRtaW4iLCJpZCI6MX0="
	GrafanaEditorUserAPIKey = "eyJrIjoiNjg1MFlVTjhodW5OYWkzVjZsOTNWcURodlA5MmhxZUQiLCJuIjoiZWRpdG9yIiwiaWQiOjF9"
	GrafanaViewerUserAPIKey = "eyJrIjoiVm00VmdEd1ZqcXVqanFFVkJiVmFCNEx0M29hR2xBSTIiLCJuIjoidmlld2VyIiwiaWQiOjF9"
	GrafanaBasicAuthUsername = "admin"
	GrafanaBasicAuthPassword = "admin"

	//Grafana default container settings
	GrafanaInternalPort = 3000
)

//Starts a grafana container that's already got the above API keys with the default test db source enabled
func StartGrafanaTestDBContainer(t *testing.T, ctx context.Context) (testcontainers.Container, string, int) {
	gRez := testcontainers.ContainerRequest{
		FromDockerfile: testcontainers.FromDockerfile{},
		Image:          "sajeevany/grafana_testdb:7.0.4",
		ExposedPorts:   []string{fmt.Sprintf("%d/tcp", GrafanaInternalPort)},
		WaitingFor:     wait.ForLog("HTTP Server Listen\" logger=http.server address=[::]:3000 protocol=http"),
	}
	grafanaC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: gRez,
		Started:          true,
	})
	if err != nil {
		t.Error(err)
	}
	grafanaIP, hErr := grafanaC.Host(ctx)
	if hErr != nil {
		t.Error(hErr)
	}
	gp := strconv.Itoa(GrafanaInternalPort)
	grafanaPort, pErr := grafanaC.MappedPort(ctx, nat.Port(gp))
	if pErr != nil {
		t.Error(pErr)
	}
	grafanaPortInt, sErr := strconv.Atoi(grafanaPort.Port())
	if sErr != nil {
		t.Error(sErr)
	}
	return grafanaC, grafanaIP, grafanaPortInt
}