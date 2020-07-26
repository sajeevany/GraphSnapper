package grafana

import (
	"context"
	"encoding/json"
	"github.com/sajeevany/graph-snapper/internal/common"
	"github.com/sajeevany/graph-snapper/internal/test"
	"github.com/sirupsen/logrus"
	"reflect"
	"testing"
	"time"
)

func TestCreateGetSnapshot(t *testing.T) {
	type args struct {
		logger    *logrus.Logger
		host      string
		port      int
		user      common.Basic
		startTime time.Time
		endTime   time.Time
		expiry    time.Time
		dashboard json.RawMessage
	}

	//Setup common requirements. In this case it's a specific grafana image with preconfigured credentials and charts
	ctx := context.Background()
	grafanaC, grafanaIP, grafanaPortInt := test.StartGrafanaTestDBContainer(t, ctx)
	defer grafanaC.Terminate(ctx)

	dashboardJson := test.GrafanaDashBJson_TXSTREZ

	tests := []struct {
		name            string
		args            args
		dashboardExists bool
	}{
		{
			name: "test0: Test create and get snapshot",
			args: args{
				logger: logrus.New(),
				host:   grafanaIP,
				port:   grafanaPortInt,
				user: common.Basic{
					Username: test.GrafanaBasicAuthUsername,
					Password: test.GrafanaBasicAuthPassword,
				},
				startTime: time.Now().Add(time.Duration(-30) * time.Minute),
				endTime:   time.Now(),
				expiry:    time.Now().AddDate(0, 0, 1),
				dashboard: json.RawMessage(dashboardJson),
			},
			dashboardExists: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			//Create snapshot
			snapshotResp, err := CreateSnapshot(tt.args.logger, tt.args.host, tt.args.port, tt.args.user, tt.args.startTime, tt.args.endTime, tt.args.expiry, tt.args.dashboard)
			if err != nil {
				t.Errorf("CreateSnapshot() error = %v", err)
				return
			}

			//Verify snapshotResp
			if snapshotResp == (CreateSnapshotResponse{}) {
				t.Errorf("Snapshot response <%+v> is empty", snapshotResp)
				return
			}

			//Get snapshot summaries
			snapshots, gsErr := GetSnapshots(tt.args.logger, tt.args.host, tt.args.port, tt.args.user)
			if gsErr != nil {
				t.Errorf("GetSnapshots() error = %v", gsErr)
				return
			}

			//Verify snapshots
			if len(snapshots) == 0 {
				t.Errorf("Grafana snapshots is empty.")
				return
			}

		})
	}
}

func TestGetSnapshotsResponse_GetSnapshot(t *testing.T) {
	tests := []struct {
		name                  string
		gsr                   GetSnapshotsResponse
		snapshotId            string
		expectSnapshotToExist bool
		expectedSnapshot      GSnapshot
	}{
		{
			//Happy path test. Expect to be able to fetch a snapshot by id from an array of snapshots
			name: "test0: Test get snapshot from a series of snapshots ",
			gsr: GetSnapshotsResponse([]GSnapshot{
				{
					ID: "0",
				},
				{
					ID: "1",
				},
			}),
			snapshotId:            "0",
			expectSnapshotToExist: true,
			expectedSnapshot: GSnapshot{
				ID: "0",
			},
		},
		{
			name: "test1: Test getSnapshot when the specified snapshot does not exist ",
			gsr: GetSnapshotsResponse([]GSnapshot{
				{
					ID: "0",
				},
				{
					ID: "1",
				},
			}),
			snapshotId:            "2",
			expectSnapshotToExist: false,
			expectedSnapshot:      GSnapshot{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			snapExists, snapshot := tt.gsr.GetSnapshot(tt.snapshotId)
			if snapExists != tt.expectSnapshotToExist {
				t.Errorf("GetSnapshot() exists check snapExists = %v, expirationInSeconds %v", snapExists, tt.expectSnapshotToExist)
			}
			if !tt.expectSnapshotToExist {
				if !reflect.DeepEqual(snapshot, tt.expectedSnapshot) {
					t.Errorf("GetSnapshot() snapshot = %v, expirationInSeconds %v", snapshot, tt.expectedSnapshot)
				}
			}
		})
	}
}

func Test_validateAndGetExpiration(t *testing.T) {

	//Define now here so that each test has a consistent point of what now is
	now := time.Now()

	type args struct {
		expiry        time.Time
		minimumExpiry time.Duration
	}
	tests := []struct {
		name                string
		args                args
		expirationInSeconds int
	}{
		{
			name: "test0: GetExpiration with expiration that exceeds minimum",
			args: args{
				expiry:        now.AddDate(0, 0, 2),
				minimumExpiry: 10 * time.Second,
			},
			expirationInSeconds: int(now.AddDate(0,0,2).Sub(now).Seconds()),
		},
		{
			name: "test1: GetExpiration with expiration that below minimum",
			args: args{
				expiry:        now.AddDate(0, 0, 1),
				minimumExpiry: 2 * 24 * time.Hour,
			},
			expirationInSeconds: int((2 * 24 * time.Hour).Seconds()),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := validateAndGetExpiration(tt.args.expiry, now, tt.args.minimumExpiry); got != tt.expirationInSeconds {
				t.Errorf("validateAndGetExpiration() = %v, expirationInSeconds %v", got, tt.expirationInSeconds)
			}
		})
	}
}