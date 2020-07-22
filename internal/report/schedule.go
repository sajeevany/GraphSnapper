package report

import (
	"github.com/sajeevany/graph-snapper/internal/common"
	"net/http"
	"time"
)

type DashboardSnapshotReport struct {
	Title            string
	Timestamp        time.Time
	GrafanaDBReports map[string]*GrafanaDashboardReport
}

func (v DashboardSnapshotReport) GetResultCode() int {
	return http.StatusOK
}

func (v DashboardSnapshotReport) ToCheckScheduleV1View() CheckV1View {
	return CheckV1View{}
}

type GrafanaDashboardReport struct {
	Title     string
	StartTime time.Time
	EndTime   time.Time
	UID       string
	Request   common.GrafanaDashBoard
	Steps     Steps
}

func (r GrafanaDashboardReport) Finalize() {
	r.EndTime = time.Now()
}

type Steps struct {
	DashboardExistsCheck  Result
	PanelSnapshot         map[string]Result
	BasicUILogin          Result
	PanelSnapshotDownload map[string]Result
	DataStorePageCreation Result
	UploadSnapshots       Result
}

type Result struct {
	Result bool
	Cause  string
}

func NewNotExecutedResult() Result{
	return Result{
		Result: false,
		Cause:  "Not executed",
	}
}