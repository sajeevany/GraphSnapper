package schedule

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/chromedp/chromedp"
	"github.com/gin-gonic/gin"
	"github.com/sajeevany/graph-snapper/internal/common"
	"github.com/sajeevany/graph-snapper/internal/grafana"
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

		report, err := recordDashboardPanels(logger, dashboard)
		snapReport.GrafanaDBReports[reqKey] = report
		if err != nil {
			logger.Errorf("Failed to perform snapshotAndUpload action for dash board <%v>. err <%v>", dashboard, err)
			continue
		}
	}

	return snapReport, nil
}

func recordDashboardPanels(logger *logrus.Logger, dashboard common.GrafanaDashBoard) (*report.GrafanaDashboardReport, error) {

	//Stub report
	dashReport := &report.GrafanaDashboardReport{
		Title:     fmt.Sprintf("Grafana dashboard <%s> snapshot panel report", dashboard.UID),
		StartTime: time.Now(),
		UID:       dashboard.UID,
		Steps: report.Steps{
			DashboardExistsCheck:  report.NewNotExecutedResult(),
			ExtractPanelID:        report.NewNotExecutedResult(),
			DashboardSnapshot:     report.NewNotExecutedResult(),
			CreateDownloadDir:     report.NewNotExecutedResult(),
			BasicUILogin:          report.NewNotExecutedResult(),
			PanelSnapshotDownload: nil,
			DataStorePageCreation: report.NewNotExecutedResult(),
			UploadSnapshots:       report.NewNotExecutedResult(),
			DeleteSnapshot:        report.NewNotExecutedResult(),
			DeleteDownloadDir:     report.NewNotExecutedResult(),
		},
	}
	defer dashReport.Finalize()

	//check if specified dashboard exists. Get the dashboard information so it can be reused to create the snapshot
	gdbExists, dashJson, dashErr := grafana.DashboardExists(logger, dashboard.UID, dashboard.Host, dashboard.Port, dashboard.User.Auth.Basic)
	if failureOccurred := setDashExistsResult(logger, dashErr, dashboard, dashReport, gdbExists); failureOccurred {
		return dashReport, dashErr
	}

	//Get panels to be screencaptured. Skip dashboard if no panels are to be captured
	panelIDs, pErr := grafana.GetPanelsIDs(dashJson, dashboard.IncludePanelsIDs, dashboard.ExcludePanelsIDs)
	if setFailedResult := setGetPanelIDsResult(logger, pErr, dashReport, panelIDs, dashboard.IncludePanelsIDs, dashboard.ExcludePanelsIDs, dashJson); setFailedResult {
		return dashReport, dashErr
	}

	//Create temporary directory to store images locally
	tmpDir, tErr := ioutil.TempDir(os.TempDir(), fmt.Sprintf("schedule-test-%s", dashboard.UID))
	if tErr != nil {
		msg := fmt.Sprintf("Failed to create temporary directory to store downloaded images <%v>", tErr)
		logger.Error(msg)
		dashReport.Steps.CreateDownloadDir = report.Result{
			Result: false,
			Cause:  msg,
		}
		return dashReport, tErr
	}
	dashReport.Steps.CreateDownloadDir = report.Result{
		Result: true,
	}
	defer func() {
		if rmErr := os.RemoveAll(tmpDir); rmErr != nil {
			dashReport.Steps.DeleteDownloadDir = report.Result{
				Result: false,
				Cause:  rmErr.Error(),
			}
		}
	}()

	//Create snapshot, screen capture and save to temporary directory
	images, cdErr := createAndDownloadSnapshotPanels(logger, dashReport, dashboard.Host, dashboard.Port, dashboard.User.Auth.Basic, dashJson, panelIDs, tmpDir)
	if cdErr != nil {
		//c&d method is responsible for updating the dashboard report
		logger.Debug("An error occurred while creating and downloading the dashboard snapshot")
		return dashReport, cdErr
	}

	logger.Info(images)

	//create page under parent page with correct file names

	//upload all attachments with name to page
	return dashReport, nil
}

func createAndDownloadSnapshotPanels(logger *logrus.Logger, dashReport *report.GrafanaDashboardReport, host string, port int, user common.Basic, dashJson json.RawMessage, panelIDs []int, storeDir string) (map[int]string, error) {

	//Opting to wrap these methods to ensure snapshot cleanup occurs
	//create dashboard snapshot
	endTime := time.Now()
	startTime := endTime.AddDate(0, 0, -1)
	expiry := time.Now().AddDate(0, 0, 0)
	gs, sErr := grafana.CreateSnapshot(logger, host, port, user, startTime, endTime, expiry, dashJson)
	if setFailedResult := setDashBoardSnapshotResult(logger, sErr, dashReport); setFailedResult {
		return nil, sErr
	}
	defer func() {
		deleteErr := grafana.DeleteSnapshot(logger, host, port, user, gs.Key)
		setDeleteSnapshotResult(logger, deleteErr, gs.Key, dashReport)
	}()

	//Download images for all panels
	files, fErr := downloadPanelImages(logger, dashReport, host, port, user, gs.Key, panelIDs, storeDir)
	if fErr != nil {
		return nil, fErr
	}

	return files, nil
}

func downloadPanelImages(logger *logrus.Logger, dashReport *report.GrafanaDashboardReport, host string, port int, user common.Basic, snapshotKey string, panelIDs []int, storeDir string) (map[int]string, error) {

	//login to dashboard
	ctxt, _ := context.WithTimeout(context.Background(), time.Minute)
	ctx, cancel := chromedp.NewContext(ctxt)
	defer cancel()

	//Login to grafana page
	loginURl := fmt.Sprintf(grafana.LoginURL, host, port)
	if err := chromedp.Run(ctx, grafana.GetLoginTasks(loginURl, user.Username, user.Password)); err != nil {
		msg := fmt.Sprintf("Unable log into grafana UI <%v>. err <%v>", loginURl, err)
		logger.Error(msg)
		dashReport.Steps.BasicUILogin = report.Result{
			Result: false,
			Cause:  msg,
		}
		return nil, err
	}

	//build url to snapshot
	urls := buildPanelUrls(host, port, snapshotKey, panelIDs)
	logger.Debugf("Going to start downloading images with panel urls <%+v>", urls)

	//screen shot each snapshot and save to local dir
	dashReport.Steps.PanelSnapshotDownload = make(map[int]report.Result, len(urls))
	images := make(map[int]string, len(panelIDs))
	for panelID, panelURL := range urls {

		//Create file name
		f, fErr := ioutil.TempFile(storeDir, "")
		if fErr != nil {
			msg := fmt.Sprintf("Unable to create temporary directory within <%v> for panel url <%v>. err <%v>", storeDir, panelURL, fErr)
			logger.Error(msg)
			dashReport.Steps.PanelSnapshotDownload[panelID] = report.Result{
				Result: false,
				Cause:  msg,
			}
			continue
		}
		logger.Debugf("Created temp file <%v> to store image from <%v>", f.Name(), panelURL)

		//Download snapshot to storage directory
		if rerr := chromedp.Run(ctx, grafana.SaveSnapshot(panelID, panelURL, f.Name())); rerr != nil {
			msg := fmt.Sprintf("Unable to open url <%v> and download snapshot. err <%v>", panelURL, rerr)
			logger.Error(msg)
			dashReport.Steps.PanelSnapshotDownload[panelID] = report.Result{
				Result: false,
				Cause:  msg,
			}
			continue
		}

		//snapshot has been saved. Update report and record
		images[panelID] = f.Name()
		dashReport.Steps.PanelSnapshotDownload[panelID] = report.Result{
			Result: true,
		}
	}

	return images, nil
}

func buildPanelUrls(host string, port int, snapshotID string, panelIDs []int) map[int]string {

	urls := make(map[int]string, len(panelIDs))

	for _, id := range panelIDs {
		urls[id] = fmt.Sprintf("http://%s:%d/dashboard/snapshot/%s?viewPanel=%d", host, port, snapshotID, id)
	}

	return urls
}

func setDeleteSnapshotResult(logger *logrus.Logger, err error, snapshotKey string, g *report.GrafanaDashboardReport) {
	if err != nil {
		msg := fmt.Sprintf("Failed to delete snapshot with key <%v>. Error <%v>", snapshotKey, err)
		logger.Error(msg)
		g.Steps.DeleteSnapshot = report.Result{
			Result: false,
			Cause:  "msg",
		}
	} else {
		g.Steps.DeleteSnapshot = report.Result{
			Result: true,
		}
	}
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
