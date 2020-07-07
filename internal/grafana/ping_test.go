package grafana

import (
	"context"
	"fmt"
	"github.com/docker/go-connections/nat"
	"github.com/sajeevany/graph-snapper/internal/common"
	"github.com/sajeevany/graph-snapper/internal/test"
	"github.com/sirupsen/logrus"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"strconv"
	"testing"
)

//TestIsValidLogin - Validates grafana login credentials check function
func TestIsValidLogin(t *testing.T) {

	//Skip test if user wants to only run regression tests
	if testing.Short(){
		t.Skip()
	}

	//Setup common requirements. In this case it's a specific grafana image with preconfigured credentials
	ctx := context.Background()
	gRez := testcontainers.ContainerRequest{
		FromDockerfile:  testcontainers.FromDockerfile{},
		Image:           "sajeevany/grafana_testdb:7.0.4",
		ExposedPorts:    []string{fmt.Sprintf("%d/tcp", test.GrafanaInternalPort)},
		WaitingFor:      wait.ForLog("HTTP Server Listen\" logger=http.server address=[::]:3000 protocol=http"),
	}
	grafanaC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: gRez,
		Started:          true,
	})
	if err != nil{
		t.Error(err)
	}
	defer grafanaC.Terminate(ctx)
	grafanaIP, hErr := grafanaC.Host(ctx)
	if hErr != nil{
		t.Error(hErr)
	}
	gp := strconv.Itoa(test.GrafanaInternalPort)
	grafanaPort, pErr := grafanaC.MappedPort(ctx, nat.Port(gp))
	if pErr != nil{
		t.Error(pErr)
	}
	grafanaPortInt, sErr := strconv.Atoi(grafanaPort.Port())
	if sErr != nil{
		t.Error(sErr)
	}

	//Scenarios
	type args struct {
		logger *logrus.Logger
		auth   common.Auth
		host   string
		port   int
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name:    "Test0 - Validate enabled admin API key",
			args:    args{
				logger: logrus.New(),
				auth:   common.Auth{
					BearerToken: common.BearerToken{
						Token: test.GrafanaAdminUserAPIKey,
					},
				},
				host:   grafanaIP,
				port:   grafanaPortInt,
			},
			want:    true,
			wantErr: false,
		},
		{
			name:    "Test1 - Validate enabled editor API key",
			args:    args{
				logger: logrus.New(),
				auth:   common.Auth{
					BearerToken: common.BearerToken{
						Token: test.GrafanaEditorUserAPIKey,
					},
				},
				host:   grafanaIP,
				port:   grafanaPortInt,
			},
			want:    true,
			wantErr: false,
		},
		{
			name:    "Test2 - Validate enabled viewer API key",
			args:    args{
				logger: logrus.New(),
				auth:   common.Auth{
					BearerToken: common.BearerToken{
						Token: test.GrafanaViewerUserAPIKey,
					},
				},
				host:   grafanaIP,
				port:   grafanaPortInt,
			},
			want:    true,
			wantErr: false,
		},
		{
			name:    "Test3 - Validate enabled basic admin credentials",
			args:    args{
				logger: logrus.New(),
				auth:   common.Auth{
					Basic: common.Basic{
						Username: test.GrafanaBasicAuthUsername,
						Password: test.GrafanaBasicAuthPassword,
					},
				},
				host:   grafanaIP,
				port:   grafanaPortInt,
			},
			want:    true,
			wantErr: false,
		},
		{
			name:    "Test4 - Validate invalid API key",
			args:    args{
				logger: logrus.New(),
				auth:   common.Auth{
					BearerToken: common.BearerToken{
						Token: "abcde",
					},
				},
				host:   grafanaIP,
				port:   grafanaPortInt,
			},
			want:    false,
			wantErr: false,
		},
		{
			name:    "Test5 - Validate invalid basic auth",
			args:    args{
				logger: logrus.New(),
				auth:   common.Auth{
					Basic: common.Basic{
						Username: "fakeUser",
						Password: "no pass",
					},
				},
				host:   grafanaIP,
				port:   grafanaPortInt,
			},
			want:    false,
			wantErr: false,
		},
	}


	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := IsValidLogin(tt.args.logger, tt.args.auth, tt.args.host, tt.args.port)
			if (err != nil) != tt.wantErr {
				t.Errorf("IsValidLogin() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("IsValidLogin() got = %v, want %v", got, tt.want)
			}
		})
	}
}