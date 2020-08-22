package schedule

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sajeevany/graph-snapper/internal/common"
	"github.com/sajeevany/graph-snapper/internal/report"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

const (
	Group         = "/schedule"
	CheckEndpoint = "/check"

	SnapshotURLFmt = "http://%s:%d/dashboard/snapshot/%s?viewPanel=%d"
)

//@Summary Check and execute schedule
//@Description Non-authenticated endpoint which checks and runs a schedule to validate connectivity and storage behaviour by the end user
//@Produce json
//@Param schedule body CheckScheduleV1 true "Check schedule"
//@Success 200 {object} report.CheckV1ReportView
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

		//Run snapshot and upload processes
		report, err := executeSchedule(logger, schedule)
		if err != nil {
			logger.WithFields(schedule.GetFields()).Errorf("Error when running executeSchedule err <%v>", err)
			return
		}

		ctx.JSON(report.GetResultCode(), report.ToCheckScheduleV1View())
	}
}

//executeSchedule - Captures dashboards and stores them to each datastore provided
func executeSchedule(logger *logrus.Logger, schedule CheckScheduleV1) (report.DashboardSnapshotReport, error) {

	snapReport := report.DashboardSnapshotReport{
		Title:            "CheckV1 schedule test",
		StartTime:        time.Now(),
		GrafanaDBReports: make(map[string]*report.GrafanaDashboardReport, len(schedule.GraphDashboards.Grafana)),
	}
	defer snapReport.Finalize()

	//create snapshots, group, and upload as a subpage to the specified datastore(s)
	for reqKey, dashboard := range schedule.GraphDashboards.Grafana {
		//wrapping and calling function to ensure defer is run per loop. TODO make this easier to read by moving it inside record function
		func() {
			report := report.NewGrafanaDashboardReport(dashboard.UID)
			defer report.Finalize()
			snapReport.GrafanaDBReports[reqKey] = report
			if recErr := recordGrafanaDashboard(logger, dashboard, schedule.DataStores, report); recErr != nil {
				logger.Errorf("Failed to record grafana dashboard with uid <%v>. error <%v>", dashboard.UID, recErr)
			}
		}()
	}

	fmt.Printf("returning report <%v>", snapReport)

	return snapReport, nil
}

func recordGrafanaDashboard(logger *logrus.Logger, dashboard common.GrafanaDashBoard, datastores DataStores, dashboardReport *report.GrafanaDashboardReport) error {
	logger.Debugf("Starting snapshot and upload for grafana dashboard <%v>", dashboard.UID)

	//Create temporary directory to store images locally. Do this here to ensure we can defer deletion.
	tmpDir, tErr := ioutil.TempDir(os.TempDir(), fmt.Sprintf("schedule-test-%s", dashboard.UID))
	if tErr != nil {
		msg := fmt.Sprintf("Failed to create temporary directory to store downloaded images <%v>", tErr)
		dashboardReport.Stages.GrafanaSnapshotStages.CreateDownloadDir = report.Result{
			Result: false,
			Cause:  msg,
		}
		logger.Error(msg)
		return tErr
	}
	dashboardReport.Stages.GrafanaSnapshotStages.CreateDownloadDir = report.Result{
		Result: true,
	}
	defer deleteDir(logger, &dashboardReport.Stages.GrafanaSnapshotStages.DeleteDownloadDir, tmpDir)

	//Create snapshot and download screenshots of requested panels
	dashboardName, images, err := captureDashboardPanels(logger, dashboardReport.Stages.GrafanaSnapshotStages, dashboard, tmpDir)
	if err != nil {
		logger.Errorf("Failed to perform captureDashboardPanels action for dash board <%v>. err <%v>", dashboard, err)
		return err
	}

	//Create/update datastore page and upload images
	uErr := setupDatastoreAndUploadDashboardPanels(logger, dashboardReport.Stages.ConfluenceStoreStages, datastores, dashboardName, dashboard.UID, images)
	if uErr != nil {
		logger.Errorf("Failed to upload images to specified datastores <%+v>. err <%v>", datastores, uErr)
		return uErr
	}
	return nil
}

func deleteDir(logger *logrus.Logger, rep *report.Result, imgDir string) {

	logger.Debugf("Starting delete directory for <%v>", imgDir)

	if rmErr := os.RemoveAll(imgDir); rmErr != nil {
		rep.Result = false
		rep.Cause = rmErr.Error()
		logger.Errorf("Delete directory operation failed for <%v> with error <%v>", imgDir, rmErr)
	} else {
		rep.Result = true
		rep.Cause = ""
	}
}
