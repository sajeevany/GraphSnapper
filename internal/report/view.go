package report

import "time"

type CheckV1ReportView struct {
	Title            string
	Timestamp        time.Time
	GrafanaDBReports map[string]GrafanaDashboardReportView
}

type GrafanaDashboardReportView struct {
}
