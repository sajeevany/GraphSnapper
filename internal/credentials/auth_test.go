package credentials

import (
	"context"
	"github.com/sajeevany/graph-snapper/internal/common"
	"github.com/sajeevany/graph-snapper/internal/test"
	"github.com/sirupsen/logrus"
	"reflect"
	"testing"
)

func Test_authGrafanaUsersIntegration(t *testing.T) {

	//Skip test if user wants to only run regression tests
	if testing.Short() {
		t.Skip()
	}

	//Setup common requirements. In this case it's a specific grafana image with preconfigured credentials
	ctx := context.Background()
	grafanaC, grafanaIP, grafanaPortInt := test.StartGrafanaTestDBContainer(t, ctx)
	defer grafanaC.Terminate(ctx)

	type args struct {
		logger *logrus.Logger
		users  []CheckUserV1
	}
	tests := []struct {
		name    string
		args    args
		want    []CheckUserResultV1
	}{
		{
			name: "test0 all invalid bearer token users",
			args: args{
				logger: logrus.New(),
				users: []CheckUserV1{
					{
						Auth: common.Auth{
							BearerToken: common.BearerToken{
								Token: "invalid token 1",
							},
						},
						Host: grafanaIP,
						Port: grafanaPortInt,
					},
					{
						Auth: common.Auth{
							BearerToken: common.BearerToken{
								Token: "invalid token 2",
							},
						},
						Host: grafanaIP,
						Port: grafanaPortInt,
					},
				},
			},
			want: []CheckUserResultV1{
				{
					Result: false,
					Cause:  "Unauthorized. Received 401",
					CheckUserV1: CheckUserV1{
						Auth: common.Auth{
							BearerToken: common.BearerToken{
								Token: "invalid token 1",
							},
						},
						Host: grafanaIP,
						Port: grafanaPortInt,
					},
				}, {
					Result: false,
					Cause:  "Unauthorized. Received 401",
					CheckUserV1: CheckUserV1{
						Auth: common.Auth{
							BearerToken: common.BearerToken{
								Token: "invalid token 2",
							},
						},
						Host: grafanaIP,
						Port: grafanaPortInt,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := authGrafanaUsers(tt.args.logger, tt.args.users)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("authGrafanaUsers() got = %v, want %v", got, tt.want)
			}
		})
	}
}
