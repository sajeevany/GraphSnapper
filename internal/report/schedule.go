package report

import (
	"fmt"
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

func (v DashboardSnapshotReport) ToCheckScheduleV1View() CheckV1ReportView {
	return CheckV1ReportView{}
}

type GrafanaDashboardReport struct {
	Title     string
	StartTime time.Time
	EndTime   time.Time
	UID       string
	Request   common.GrafanaDashBoard
	Steps     *Steps
}

func NewGrafanaDashboardReport(uid string) *GrafanaDashboardReport {
	//Stub report
	return &GrafanaDashboardReport{
		Title:     fmt.Sprintf("Grafana dashboard <%s> snapshot panel report", uid),
		StartTime: time.Now(),
		UID:       uid,
		Steps: &Steps{
			GrafanaSnapshotSteps: &GrafanaDBSnapshotStages{
				DashboardExistsCheck:  NewNotExecutedResult(),
				ExtractPanelID:        NewNotExecutedResult(),
				DashboardSnapshot:     NewNotExecutedResult(),
				CreateDownloadDir:     NewNotExecutedResult(),
				BasicUILogin:          NewNotExecutedResult(),
				PanelSnapshotDownload: nil,
				DeleteSnapshot:        NewNotExecutedResult(),
				DeleteDownloadDir:     NewNotExecutedResult(),
			},
			ConfluenceStoreStages: nil,
		},
	}
}

func (r GrafanaDashboardReport) Finalize() {
	r.EndTime = time.Now()
}

type Steps struct {
	GrafanaSnapshotSteps  *GrafanaDBSnapshotStages
	ConfluenceStoreStages map[string]ConfluenceStoreStages
}

type GrafanaDBSnapshotStages struct {
	DashboardExistsCheck  Result
	ExtractPanelID        Result
	DashboardSnapshot     Result
	CreateDownloadDir     Result
	BasicUILogin          Result
	PanelSnapshotDownload map[int]PanelDownload

	//Cleanup stages
	DeleteSnapshot    Result
	DeleteDownloadDir Result
}

type PanelDownload struct {
	CreateTempFile          Result
	DownloadPanelScreenshot Result
}

type ConfluenceStoreStages struct {
	ParentPageExistsCheck   Result
	CreateMissingParentPage Result
	DataStorePageCreation   Result
	UploadSnapshots         Result
}

type Result struct {
	Result bool
	Cause  string
}

func NewNotExecutedResult() Result {
	return Result{
		Result: false,
		Cause:  "Not executed",
	}
}
