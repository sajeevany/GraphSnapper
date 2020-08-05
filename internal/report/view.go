package report

import "time"

type CheckV1ReportView struct {
	Title            string
	StartTime        time.Time
	EndTime          time.Time
	Duration         time.Duration
	GrafanaDBReports map[string]GrafanaDashboardReportView
}

type GrafanaDashboardReportView struct {
	Title     string
	StartTime time.Time
	EndTime   time.Time
	UID       string
	Stages    GrafanaStagesView
}

type GrafanaStagesView struct {
	GrafanaSnapshotStages GrafanaSnapshotStagesView
	ConfluenceStoreStages map[string]ConfluenceStoreStagesView
}

type GrafanaSnapshotStagesView struct {
	DashboardExistsCheck  Result
	ExtractPanelID        Result
	DashboardSnapshot     Result
	CreateDownloadDir     Result
	BasicUILogin          Result
	PanelSnapshotDownload map[int]PanelDownloadView

	//Cleanup stages
	DeleteSnapshot    Result
	DeleteDownloadDir Result
}

type PanelDownloadView struct {
	CreateTempFile          Result
	DownloadPanelScreenshot Result
}

type ConfluenceStoreStagesView struct {
	ParentPageExistsCheck   Result
	CreateMissingParentPage Result
	DataStorePageCreation   Result
	UploadSnapshots         Result
}
