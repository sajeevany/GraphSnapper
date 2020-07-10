package grafana

import (
	"context"
	"github.com/sajeevany/graph-snapper/internal/common"
	"github.com/sajeevany/graph-snapper/internal/test"
	"github.com/sirupsen/logrus"
	"testing"
)

//TestIsValidLogin - Validates grafana login credentials check function
func TestIsValidLogin(t *testing.T) {

	//Skip test if user wants to only run regression tests
	if testing.Short() {
		t.Skip()
	}

	//Setup common requirements. In this case it's a specific grafana image with preconfigured credentials
	ctx := context.Background()
	grafanaC, grafanaIP, grafanaPortInt := test.StartGrafanaTestDBContainer(t, ctx)
	defer grafanaC.Terminate(ctx)

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
			name: "Test0 - Validate enabled admin API key",
			args: args{
				logger: logrus.New(),
				auth: common.Auth{
					BearerToken: common.BearerToken{
						Token: test.GrafanaAdminUserAPIKey,
					},
				},
				host: grafanaIP,
				port: grafanaPortInt,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Test1 - Validate enabled editor API key",
			args: args{
				logger: logrus.New(),
				auth: common.Auth{
					BearerToken: common.BearerToken{
						Token: test.GrafanaEditorUserAPIKey,
					},
				},
				host: grafanaIP,
				port: grafanaPortInt,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Test2 - Validate enabled viewer API key",
			args: args{
				logger: logrus.New(),
				auth: common.Auth{
					BearerToken: common.BearerToken{
						Token: test.GrafanaViewerUserAPIKey,
					},
				},
				host: grafanaIP,
				port: grafanaPortInt,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Test3 - Validate enabled basic admin credentials",
			args: args{
				logger: logrus.New(),
				auth: common.Auth{
					Basic: common.Basic{
						Username: test.GrafanaBasicAuthUsername,
						Password: test.GrafanaBasicAuthPassword,
					},
				},
				host: grafanaIP,
				port: grafanaPortInt,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Test4 - Validate invalid API key",
			args: args{
				logger: logrus.New(),
				auth: common.Auth{
					BearerToken: common.BearerToken{
						Token: "abcde",
					},
				},
				host: grafanaIP,
				port: grafanaPortInt,
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "Test5 - Validate invalid basic auth",
			args: args{
				logger: logrus.New(),
				auth: common.Auth{
					Basic: common.Basic{
						Username: "fakeUser",
						Password: "no pass",
					},
				},
				host: grafanaIP,
				port: grafanaPortInt,
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
