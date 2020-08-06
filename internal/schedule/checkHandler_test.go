package schedule

import (
	"fmt"
	"github.com/sajeevany/graph-snapper/internal/grafana"
	"reflect"
	"testing"
)

func Test_setPanelSnapshotUrls(t *testing.T) {
	type args struct {
		host       string
		port       int
		snapshotID string
		panelIDs   []grafana.PanelDescriptor
	}
	tests := []struct {
		name string
		args args
		want []grafana.PanelDescriptor
	}{
		{
			name: "test0: Happy path",
			args: args{
				host:       "localhost",
				port:       9000,
				snapshotID: "abcdedfg",
				panelIDs: []grafana.PanelDescriptor{
					{
						Title: "Panel1",
						ID:    1,
					},
					{
						Title: "Panel2",
						ID:    2,
					},
				},
			},
			want: []grafana.PanelDescriptor{
				{
					Title:       "Panel1",
					ID:          1,
					SnapshotURL: fmt.Sprintf(SnapshotURLFmt, "localhost", 9000, "abcdedfg", 1),
				},
				{
					Title:       "Panel2",
					ID:          2,
					SnapshotURL: fmt.Sprintf(SnapshotURLFmt, "localhost", 9000, "abcdedfg", 2),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			panels := setPanelSnapshotUrls(tt.args.host, tt.args.port, tt.args.snapshotID, tt.args.panelIDs)

			if !reflect.DeepEqual(panels, tt.want) {
				t.Errorf("setPanelSnapshotUrls() actual = %v, want %v", tt.args.panelIDs, tt.want)
			}
		})
	}
}
