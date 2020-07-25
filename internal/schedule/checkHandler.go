package schedule

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sajeevany/graph-snapper/internal/common"
	"github.com/sajeevany/graph-snapper/internal/grafana"
	"github.com/sajeevany/graph-snapper/internal/report"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

const (
	Group         = "/schedule"
	CheckEndpoint = "/check"
)

//@Summary Check and execute schedule
//@Description Non-authenticated endpoint which checks and runs a schedule to validate connectivity and storage behaviour by the end user
//@Produce json
//@Param schedule body CheckScheduleV1 true "Check schedule"
//@Success 200 {object} report.CheckV1View
//@Fail 400 {object} gin.H
//@Fail 500 {object} gin.H
//@Router /schedule/check [post]
//@Tags schedule
func CheckV1(logger *logrus.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		logger.Debug("Starting schedule check (v1)")

		//Bind schedule object
		var schedule CheckScheduleV1
		if bErr := ctx.BindJSON(&schedule); bErr != nil {
			msg := fmt.Sprintf("Unable to bind request body to schedule object %v", bErr)
			logger.Errorf(msg)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": msg})
			return
		}

		//Quick check if all required attributes are present
		if isValid, err := schedule.IsValid(); !isValid {
			logger.Debug("Invalid schedule provided")
			ctx.JSON(http.StatusBadRequest, err)
			return
		}

		//Run snapshot and upload process
		report, err := snapshotAndUpload(logger, schedule)
		if err != nil {
			logger.WithFields(schedule.GetFields()).Errorf("Error when running snapshotAndUpload err <%v>", err)
			return
		}

		ctx.JSON(report.GetResultCode(), report.ToCheckScheduleV1View())

	}
}

//snapshotAndUpload -
func snapshotAndUpload(logger *logrus.Logger, schedule CheckScheduleV1) (report.DashboardSnapshotReport, error) {

	snapReport := report.DashboardSnapshotReport{
		Title:            "CheckV1 schedule test",
		Timestamp:        time.Now(),
		GrafanaDBReports: make(map[string]*report.GrafanaDashboardReport, len(schedule.GraphDashboards.Grafana)),
	}

	//create snapshots, group, and upload as a subpage to the specified datastore(s)
	for reqKey, dashboard := range schedule.GraphDashboards.Grafana {

		logger.Debugf("Starting snapshot and upload for grafana dashboard <%v>", reqKey)

		dashReport := report.GrafanaDashboardReport{
			Title:     fmt.Sprintf("Grafana dashboard <%s> snapshot panel report", dashboard.UID),
			StartTime: time.Now(),
			UID:       dashboard.UID,
			Steps: report.Steps{
				DashboardExistsCheck:  report.NewNotExecutedResult(),
				ExtractPanelID:        report.NewNotExecutedResult(),
				DashboardSnapshot:     report.NewNotExecutedResult(),
				BasicUILogin:          report.NewNotExecutedResult(),
				PanelSnapshotDownload: nil,
				DataStorePageCreation: report.NewNotExecutedResult(),
				UploadSnapshots:       report.NewNotExecutedResult(),
			},
		}
		snapReport.GrafanaDBReports[reqKey] = &dashReport

		//check if specified dashboard exists. Get the dashboard information so it can be reused to create the snapshot
		gdbExists, dashJson, dashErr := grafana.DashboardExists(logger, dashboard.UID, dashboard.Host, dashboard.Port, dashboard.User.Auth.Basic)
		if setFailedResult := setDashExistsResult(logger, dashErr, dashboard, &dashReport, gdbExists); setFailedResult {
			continue
		}

		//Get panels to be screencaptured. Skip dashboard if no panels are to be captured
		panelIDs, pErr := grafana.GetPanelsIDs(dashJson, dashboard.IncludePanelsIDs, dashboard.ExcludePanelsIDs)
		if setFailedResult := setGetPanelIDsResult(logger, pErr, &dashReport, panelIDs, dashboard.IncludePanelsIDs, dashboard.ExcludePanelsIDs, dashJson); setFailedResult {
			continue
		}

		//create dashboard snapshot
		endTime := time.Now()
		startTime := endTime.AddDate(0, 0, -1)
		expiry := time.Now().AddDate(0, 0, 0)
		gs, sErr := grafana.CreateSnapshot(logger, dashboard.Host, dashboard.Port, dashboard.User.Auth.Basic, startTime, endTime, expiry, dashJson)
		if setFailedResult := setDashBoardSnapshotResult(logger, sErr, &dashReport); setFailedResult {
			continue
		}

		//run chromedb login

		//build url to snapshot

		//screen shot each snapshot and save to local dir

		//create page under parent page with correct file names

		//upload all attachments with name to page

		//Close off report
		dashReport.Finalize()

	}

	return snapReport, nil
}

func setDashBoardSnapshotResult(logger *logrus.Logger, err error, g *report.GrafanaDashboardReport) bool {

	if err != nil {
		logger.Errorf("Unable to create snapshot for dashboard <%v>. err <%v>", g.UID, err)
		g.Steps.DashboardSnapshot = report.Result{
			Result: false,
			Cause:  err.Error(),
		}
		return true
	}

	return false
}

func setGetPanelIDsResult(logger *logrus.Logger, pErr error, dashReport *report.GrafanaDashboardReport, ids, includeIDs, excludeIDs []int, dashJson json.RawMessage) bool {

	if pErr != nil {
		logger.Errorf("Unable to parse grafana dashboard API response. <%v>", pErr)
		dashReport.Steps.ExtractPanelID = report.Result{
			Result: false,
			Cause:  pErr.Error(),
		}
		return true
	}

	if len(ids) == 0 {
		//No panels left to be recorded
		msg := fmt.Sprintf("No panels ids remaining after applying inclusion <%v> and exclusion <%v> lists to dashboard result <%v>", includeIDs, excludeIDs, dashJson)
		logger.Info(msg)
		dashReport.Steps.ExtractPanelID = report.Result{
			Result: false,
			Cause:  msg,
		}
		return true
	}

	dashReport.Steps.ExtractPanelID = report.Result{
		Result: true,
	}
	return false
}

func setDashExistsResult(logger *logrus.Logger, dashErr error, dashboard common.GrafanaDashBoard, dashReport *report.GrafanaDashboardReport, gdbExists bool) bool {
	if dashErr != nil {
		logger.Errorf("Internal error <%v> when checking if dashboard <%v> exists at host <%v> port <%v>", dashErr, dashboard.UID, dashboard.Host, dashboard.Port)
		dashReport.Steps.DashboardExistsCheck = report.Result{
			Result: false,
			Cause:  dashErr.Error(),
		}
		return true
	}
	if !gdbExists {
		msg := fmt.Sprintf("Dashboard <%v> does not exist at host <%v> port <%v>", dashboard.UID, dashboard.Host, dashboard.Port)
		logger.Debug(msg)
		dashReport.Steps.DashboardExistsCheck = report.Result{
			Result: false,
			Cause:  msg,
		}
		return true
	}
	dashReport.Steps.DashboardExistsCheck = report.Result{
		Result: true,
	}
	return false
}
