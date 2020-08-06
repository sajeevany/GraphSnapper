package report

import (
	"fmt"
	"github.com/sajeevany/graph-snapper/internal/common"
	"net/http"
	"time"
)

type DashboardSnapshotReport struct {
	Title            string
	StartTime        time.Time
	EndTime          time.Time
	GrafanaDBReports map[string]*GrafanaDashboardReport
}

func (v DashboardSnapshotReport) GetResultCode() int {
	return http.StatusOK
}

func (v DashboardSnapshotReport) ToCheckScheduleV1View() CheckV1ReportView {
	return CheckV1ReportView{
		Title:     v.Title,
		StartTime: v.StartTime,
		EndTime:   v.EndTime,
		Duration:  v.EndTime.Sub(v.EndTime),
		GrafanaDBReports: func() map[string]GrafanaDashboardReportView {
			view := make(map[string]GrafanaDashboardReportView, len(v.GrafanaDBReports))
			for key, rep := range v.GrafanaDBReports {
				view[key] = rep.toGrafanaDashboardReportView()
			}
			return view
		}(),
	}
}

func (v *DashboardSnapshotReport) Finalize() {
	v.EndTime = time.Now()
}

type GrafanaDashboardReport struct {
	Title     string
	StartTime time.Time
	EndTime   time.Time
	UID       string
	Request   common.GrafanaDashBoard
	Stages    *Stages
}

func (g *GrafanaDashboardReport) toGrafanaDashboardReportView() GrafanaDashboardReportView {
	return GrafanaDashboardReportView{
		Title:     g.Title,
		StartTime: g.StartTime,
		EndTime:   g.EndTime,
		UID:       g.UID,
		Stages:    g.Stages.toStagesView(),
	}
}

func NewGrafanaDashboardReport(uid string) *GrafanaDashboardReport {
	//Stub report
	return &GrafanaDashboardReport{
		Title:     fmt.Sprintf("Grafana dashboard <%s> snapshot panel report", uid),
		StartTime: time.Now(),
		UID:       uid,
		Stages: &Stages{
			GrafanaSnapshotStages: &GrafanaDBSnapshotStages{
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

func (r *GrafanaDashboardReport) Finalize() {
	r.EndTime = time.Now()
}

type Stages struct {
	GrafanaSnapshotStages *GrafanaDBSnapshotStages
	ConfluenceStoreStages map[string]ConfluenceStoreStages
}

func (s Stages) toStagesView() GrafanaStagesView {
	return GrafanaStagesView{
		GrafanaSnapshotStages: s.GrafanaSnapshotStages.toGrafanaSnapshotStagesView(),
		ConfluenceStoreStages: func() map[string]ConfluenceStoreStagesView {

			m := make(map[string]ConfluenceStoreStagesView, len(s.ConfluenceStoreStages))
			for i, v := range s.ConfluenceStoreStages {
				m[i] = v.ConfluenceStoreStagesView()
			}

			return m
		}(),
	}
}

type GrafanaDBSnapshotStages struct {
	DashboardExistsCheck  Result
	ExtractPanelID        Result
	DashboardSnapshot     Result
	CreateDownloadDir     Result
	BasicUILogin          Result
	PanelSnapshotDownload map[int]*PanelDownload

	//Cleanup stages
	DeleteSnapshot    Result
	DeleteDownloadDir Result
}

func (gs GrafanaDBSnapshotStages) toGrafanaSnapshotStagesView() GrafanaSnapshotStagesView {
	return GrafanaSnapshotStagesView{
		DashboardExistsCheck:  gs.DashboardExistsCheck,
		ExtractPanelID:        gs.ExtractPanelID,
		DashboardSnapshot:     gs.DashboardSnapshot,
		CreateDownloadDir:     gs.CreateDownloadDir,
		BasicUILogin:          gs.BasicUILogin,
		PanelSnapshotDownload: ToPanelDownloadViewMap(gs.PanelSnapshotDownload),
		DeleteSnapshot:        gs.DeleteSnapshot,
		DeleteDownloadDir:     gs.DeleteDownloadDir,
	}
}

type PanelDownload struct {
	CreateTempFile          Result
	DownloadPanelScreenshot Result
}

func ToPanelDownloadViewMap(pd map[int]*PanelDownload) map[int]PanelDownloadView {
	m := make(map[int]PanelDownloadView, len(pd))
	for i, v := range pd {
		m[i] = v.ToPanelDownloadView()
	}
	return m
}

func (pd PanelDownload) ToPanelDownloadView() PanelDownloadView {
	return PanelDownloadView{
		CreateTempFile:          pd.CreateTempFile,
		DownloadPanelScreenshot: pd.DownloadPanelScreenshot,
	}
}

type ConfluenceStoreStages struct {
	ParentPageExistsCheck   Result
	CreateMissingParentPage Result
	DataStorePageCreation   Result
	UploadSnapshots         Result
}

func (c ConfluenceStoreStages) ConfluenceStoreStagesView() ConfluenceStoreStagesView {
	return ConfluenceStoreStagesView{
		ParentPageExistsCheck:   c.ParentPageExistsCheck,
		CreateMissingParentPage: c.CreateMissingParentPage,
		DataStorePageCreation:   c.DataStorePageCreation,
		UploadSnapshots:         c.UploadSnapshots,
	}
}

func NewConfluenceStoreStages() ConfluenceStoreStages {
	return ConfluenceStoreStages{
		ParentPageExistsCheck:   NewNotExecutedResult(),
		CreateMissingParentPage: NewNotExecutedResult(),
		DataStorePageCreation:   NewNotExecutedResult(),
		UploadSnapshots:         NewNotExecutedResult(),
	}
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
