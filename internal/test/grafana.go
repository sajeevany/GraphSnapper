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
	GrafanaAdminUserAPIKey   = "eyJrIjoiUDI4NFg3WFpoNDJXRnJ4RTBZeDMwb2U0WnA2QTBkRmoiLCJuIjoiYWRtaW4iLCJpZCI6MX0="
	GrafanaEditorUserAPIKey  = "eyJrIjoiNjg1MFlVTjhodW5OYWkzVjZsOTNWcURodlA5MmhxZUQiLCJuIjoiZWRpdG9yIiwiaWQiOjF9"
	GrafanaViewerUserAPIKey  = "eyJrIjoiVm00VmdEd1ZqcXVqanFFVkJiVmFCNEx0M29hR2xBSTIiLCJuIjoidmlld2VyIiwiaWQiOjF9"
	GrafanaBasicAuthUsername = "admin"
	GrafanaBasicAuthPassword = "admin2"

	//Grafana default container settings
	GrafanaInternalPort = 3000

	//dashboard json
	GrafanaDashBJson_TXSTREZ = "{ \"__requires\": [ { \"id\": \"grafana\", \"name\": \"Grafana\", \"type\": \"grafana\", \"version\": \"6.3.0-pre\" }, { \"id\": \"graph\", \"name\": \"Graph\", \"type\": \"panel\", \"version\": \"\" }, { \"id\": \"graph\", \"name\": \"React Graph\", \"type\": \"panel\", \"version\": \"\" }, { \"id\": \"testdata\", \"name\": \"TestData DB\", \"type\": \"datasource\", \"version\": \"1.0.0\" } ], \"annotations\": { \"list\": [ { \"builtIn\": 1, \"datasource\": \"-- Grafana --\", \"enable\": true, \"hide\": true, \"iconColor\": \"rgba(0, 211, 255, 1)\", \"name\": \"Annotations & Alerts\", \"type\": \"dashboard\" } ] }, \"editable\": true, \"gnetId\": null, \"graphTooltip\": 0, \"id\": 1, \"links\": [], \"panels\": [ { \"aliasColors\": {}, \"bars\": false, \"dashLength\": 10, \"dashes\": false, \"fill\": 1, \"gridPos\": { \"h\": 7, \"w\": 24, \"x\": 0, \"y\": 0 }, \"id\": 4, \"legend\": { \"alignAsTable\": true, \"avg\": false, \"current\": true, \"max\": false, \"min\": false, \"rightSide\": true, \"show\": true, \"total\": false, \"values\": true }, \"lines\": true, \"linewidth\": 1, \"links\": [], \"nullPointMode\": \"null\", \"percentage\": false, \"pointradius\": 2, \"points\": false, \"renderer\": \"flot\", \"seriesOverrides\": [], \"spaceLength\": 10, \"stack\": false, \"steppedLine\": false, \"targets\": [ { \"refId\": \"A\", \"scenarioId\": \"streaming_client\", \"stream\": { \"noise\": 2.2, \"speed\": 100, \"spread\": 3.5, \"type\": \"signal\" }, \"stringInput\": \"\" } ], \"thresholds\": [], \"timeFrom\": null, \"timeRegions\": [], \"timeShift\": null, \"title\": \"Angular Graph\", \"tooltip\": { \"shared\": true, \"sort\": 0, \"value_type\": \"individual\" }, \"type\": \"graph\", \"xaxis\": { \"buckets\": null, \"mode\": \"time\", \"name\": null, \"show\": true, \"values\": [] }, \"yaxes\": [ { \"decimals\": 2, \"format\": \"short\", \"label\": null, \"logBase\": 1, \"max\": null, \"min\": null, \"show\": true }, { \"format\": \"short\", \"label\": null, \"logBase\": 1, \"max\": null, \"min\": null, \"show\": true } ], \"yaxis\": { \"align\": false, \"alignLevel\": null } }, { \"datasource\": \"TestDataDB\", \"description\": \"\", \"gridPos\": { \"h\": 6, \"w\": 24, \"x\": 0, \"y\": 7 }, \"id\": 2, \"links\": [], \"options\": { \"graph\": { \"showBars\": false, \"showLines\": true, \"showPoints\": false }, \"legend\": { \"asTable\": true, \"decimals\": 2, \"isVisible\": true, \"placement\": \"right\", \"stats\": [ \"last\" ] }, \"series\": {} }, \"targets\": [ { \"refId\": \"A\", \"scenarioId\": \"streaming_client\", \"stream\": { \"noise\": 10, \"speed\": 100, \"spread\": 20, \"type\": \"signal\" }, \"stringInput\": \"\" } ], \"timeFrom\": null, \"timeShift\": null, \"title\": \"Simple dummy streaming example\", \"type\": \"graph\" } ], \"schemaVersion\": 18, \"style\": \"dark\", \"tags\": [], \"templating\": { \"list\": [] }, \"time\": { \"from\": \"now-1m\", \"to\": \"now\" }, \"timepicker\": { \"refresh_intervals\": [ \"5s\", \"10s\", \"30s\", \"1m\", \"5m\", \"15m\", \"30m\", \"1h\", \"2h\", \"1d\" ], \"time_options\": [ \"5m\", \"15m\", \"1h\", \"6h\", \"12h\", \"24h\", \"2d\", \"7d\", \"30d\" ] }, \"timezone\": \"\", \"title\": \"Simple Streaming Example\", \"uid\": \"TXSTREZ\", \"version\": 1 }"
)

//Starts a grafana container that's already got the above API keys with the default test db source enabled
func StartGrafanaTestDBContainer(t *testing.T, ctx context.Context) (testcontainers.Container, string, int) {
	gRez := testcontainers.ContainerRequest{
		FromDockerfile: testcontainers.FromDockerfile{},
		Image:          "sajeevany/grafana_testdb:7.0.4.1",
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
