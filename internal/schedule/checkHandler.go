package schedule

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/chromedp/chromedp"
	"github.com/gin-gonic/gin"
	"github.com/sajeevany/graph-snapper/internal/common"
	"github.com/sajeevany/graph-snapper/internal/confluence"
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

	//Create snapshto and download screenshots of requested panels
	images, err := captureDashboardPanels(logger, dashboardReport.Stages.GrafanaSnapshotStages, dashboard, tmpDir)
	if err != nil {
		logger.Errorf("Failed to perform captureDashboardPanels action for dash board <%v>. err <%v>", dashboard, err)
		return err
	}

	//Create/update datastore page and upload images
	uErr := setupDatastoreAndUpload(logger, dashboardReport.Stages.ConfluenceStoreStages, datastores, images)
	if uErr != nil {
		logger.Errorf("Failed to upload images to specified datastores <%+v>. err <%v>", datastores, uErr)
		return uErr
	}
	return nil
}

func setupDatastoreAndUpload(logger *logrus.Logger, rep map[string]report.ConfluenceStoreStages, stores DataStores, images []grafana.DownloadedPanelDesc) error {

	logger.Debug("Starting image upload to datastores()")

	//Create and upload images to confluence page
	rep = make(map[string]report.ConfluenceStoreStages, len(stores.ConfluencePages))
	for _, parent := range stores.ConfluencePages {

		reprt := report.NewConfluenceStoreStages()
		rep[parent.ParentPageID] = reprt

		//Check if the parent confluence page exists. If not create it
		exists, eErr := confluence.DoesPageExistByID(logger, parent.ParentPageID, parent.Host, parent.Port, common.Auth{
			Basic: parent.User,
		})
		if failed := setParentCheckResult(logger, eErr, parent, &reprt.ParentPageExistsCheck, exists); failed {
			continue
		}

		//
		//now := time.Now().Format(time.RFC1123)
		//pageName := fmt.Sprintf("%s_%s", dashboardUID, strings.Replace(now, " ", "", -1))
		//
		////Create page
		//pageID, pErr := confluence.CreatePage(logger, pageName, parent.SpaceKey, parent.ParentPageID, parent.User, images)
		//if pErr != nil {
		//	//rep.Stages.
		//}

		//

	}

	return nil
}

//setParentCheckResult - sets result of parent exists check. Returns true if an error occurred or the parent page doesn't exist. Use to continue/skip any further operations.
func setParentCheckResult(logger *logrus.Logger, eErr error, parent ParentConfluencePage, pExistsResult *report.Result, exists bool) bool {

	if eErr != nil {
		msg := fmt.Sprintf("Error checking if confluence page with source id <%v> does not exist. <%v>", parent.ParentPageID, eErr.Error())
		logger.Errorf(msg)
		pExistsResult.Result = false
		pExistsResult.Cause = msg

		return true
	}
	if !exists {
		pExistsResult.Result = false
		pExistsResult.Cause = "Page does not exist"
		logger.Debugf("Page <%v> does not exist", parent.ParentPageID)

		return true
	}

	//Set passed result
	pExistsResult.Result = true
	pExistsResult.Cause = ""

	return false
}

//Snapshots a grafana dashboard and stores the images in a temporary folder. Returns a map of image descriptors to image locations and image storage directory.
func captureDashboardPanels(logger *logrus.Logger, dashboardStages *report.GrafanaDBSnapshotStages, dashboard common.GrafanaDashBoard, downloadDir string) ([]grafana.DownloadedPanelDesc, error) {

	//check if specified dashboard exists. Get the dashboard information so it can be reused to create the snapshot
	gdbExists, dashJson, dashErr := grafana.DashboardExists(logger, dashboard.UID, dashboard.Host, dashboard.Port, dashboard.User.Auth.Basic)
	if failureOccurred := setDashExistsResult(logger, dashErr, dashboard, &dashboardStages.DashboardExistsCheck, gdbExists); failureOccurred {
		return nil, dashErr
	}

	//Get panels to be screencaptured. Skip dashboard if no panels are to be captured
	panelDesc, pErr := grafana.GetPanelsDescriptors(logger, dashJson, dashboard.IncludePanelsIDs, dashboard.ExcludePanelsIDs)
	if setFailedResult := setGetPanelIDsResult(logger, pErr, &dashboardStages.ExtractPanelID, panelDesc, dashboard.IncludePanelsIDs, dashboard.ExcludePanelsIDs, dashJson); setFailedResult {
		return nil, dashErr
	}

	//Create snapshot, screen capture and save to temporary directory
	images, cdErr := createAndDownloadSnapshotPanels(logger, dashboardStages, dashboard.Host, dashboard.Port, dashboard.User.Auth.Basic, dashboard.UID, dashJson, panelDesc, downloadDir)
	if cdErr != nil {
		//c&d method is responsible for updating the dashboard report
		logger.Debugf("An error occurred while creating and downloading the dashboard snapshot <%v>", cdErr)
		return nil, cdErr
	}
	logger.Debugf("Downloaded images as <%+v>", images)

	//upload all attachments with name to page
	return images, nil
}

//createAndDownloadSnapshotPanels - create grafana dashboard snapshot and download the specified panels
func createAndDownloadSnapshotPanels(logger *logrus.Logger, dashStages *report.GrafanaDBSnapshotStages, host string, port int, user common.Basic, dashUID string, dashJson json.RawMessage, panelDesc []grafana.PanelDescriptor, storeDir string) ([]grafana.DownloadedPanelDesc, error) {

	//Opting to wrap these methods to ensure snapshot cleanup occurs
	//create dashboard snapshot. Get snapshot response because snapshot creation requires the database description
	endTime := time.Now()
	startTime := endTime.AddDate(0, 0, -1)
	expiry := time.Now().AddDate(0, 0, 0)
	gs, sErr := grafana.CreateSnapshot(logger, host, port, user, startTime, endTime, expiry, dashJson)
	if setFailedResult := setDashBoardSnapshotResult(logger, sErr, dashUID, &dashStages.DashboardSnapshot); setFailedResult {
		return nil, sErr
	}
	defer func() {
		deleteErr := grafana.DeleteSnapshot(logger, host, port, user, gs.Key)
		setDeleteSnapshotResult(logger, deleteErr, gs.Key, &dashStages.DeleteSnapshot)
	}()

	//Download images for all panels
	files, fErr := downloadPanelImages(logger, dashStages, host, port, user, gs.Key, panelDesc, storeDir)
	if fErr != nil {
		return nil, fErr
	}

	return files, nil
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

func downloadPanelImages(logger *logrus.Logger, snapshotStages *report.GrafanaDBSnapshotStages, host string, port int, user common.Basic, snapshotKey string, panels []grafana.PanelDescriptor, storeDir string) ([]grafana.DownloadedPanelDesc, error) {

	logger.Debug("Started downloadPanelImages()")
	defer logger.Debug("Completed downloadPanelImages()")

	//login to dashboard
	ctxt, _ := context.WithTimeout(context.Background(), time.Minute)
	ctx, cancel := chromedp.NewContext(ctxt)
	defer cancel()

	//Login to grafana page
	loginURl := fmt.Sprintf(grafana.LoginURL, host, port)
	if err := chromedp.Run(ctx, grafana.GetLoginTasks(loginURl, user.Username, user.Password)); err != nil {
		msg := fmt.Sprintf("Unable log into grafana UI <%v>. err <%v>", loginURl, err)
		logger.Error(msg)
		snapshotStages.BasicUILogin = report.Result{
			Result: false,
			Cause:  msg,
		}
		return nil, err
	} else {
		snapshotStages.BasicUILogin = report.Result{
			Result: true,
			Cause:  "",
		}
	}

	//build url to snapshot
	panels = setPanelSnapshotUrls(host, port, snapshotKey, panels)
	logger.Debugf("Going to start downloading images with panel urls <%+v>", panels)

	//screen shot each snapshot and save to local dir
	snapshotStages.PanelSnapshotDownload = make(map[int]*report.PanelDownload, len(panels))
	images := make([]grafana.DownloadedPanelDesc, len(panels))
	for idx, panel := range panels {

		panelDownloadResult := &report.PanelDownload{}
		snapshotStages.PanelSnapshotDownload[panel.ID] = panelDownloadResult

		//Create temp file and download panel to it
		fName, pErr := downloadToTempFile(logger, panelDownloadResult, ctx, panel, storeDir)
		if pErr != nil {
			logger.Debug("Error detected when downloading panel <%v> to temporary directory. err <%v>", panel, pErr)
			continue
		}

		//Update slice of downloaded panel images to be referenced for upload later
		images[idx] = grafana.DownloadedPanelDesc{
			PanelDescriptor: panel,
			DownloadDir:     fName,
		}

	}

	return images, nil
}

func downloadToTempFile(logger *logrus.Logger, result *report.PanelDownload, ctx context.Context, panel grafana.PanelDescriptor, dir string) (string, error) {
	//Create file name
	f, fErr := ioutil.TempFile(dir, "")
	if setCreateTempFileResult(logger, fErr, dir, panel, &result.CreateTempFile, f) {
		return "", fErr
	}

	//Download snapshot to storage directory
	rErr := chromedp.Run(ctx, grafana.SaveSnapshot(panel, f.Name()))
	if errDetected := setSaveSnapshotResult(logger, rErr, &result.DownloadPanelScreenshot, panel.SnapshotURL); errDetected {
		return "", rErr
	}

	return f.Name(), nil
}

//Updates result if an error exists. Returns true if error is non-nil
func setSaveSnapshotResult(logger *logrus.Logger, err error, r *report.Result, snapshotURL string) bool {
	if err != nil {
		msg := fmt.Sprintf("Unable to open url <%v> and download snapshot. err <%v>", snapshotURL, err)
		logger.Error(msg)
		r.Result = false
		r.Cause = msg

		return true
	} else {
		r.Result = true
		r.Cause = ""

		return false
	}
}

func setCreateTempFileResult(logger *logrus.Logger, fErr error, storeDir string, panel grafana.PanelDescriptor, res *report.Result, f *os.File) bool {
	if fErr != nil {
		msg := fmt.Sprintf("Unable to create temporary directory within <%v> for panel url <%v>. err <%v>", storeDir, panel.SnapshotURL, fErr)
		logger.Error(msg)
		res.Result = false
		res.Cause = msg

		return true
	} else {
		//Temp file created. Updated result
		res.Result = true
		res.Cause = ""
		logger.Debugf("Created temp file <%v> to store image from <%v>", f.Name(), panel.SnapshotURL)

		return false
	}
}

func setPanelSnapshotUrls(host string, port int, snapshotID string, panelIDs []grafana.PanelDescriptor) []grafana.PanelDescriptor {
	for idx, panel := range panelIDs {
		panelIDs[idx].SnapshotURL = fmt.Sprintf(SnapshotURLFmt, host, port, snapshotID, panel.ID)
	}
	return panelIDs
}

func setDeleteSnapshotResult(logger *logrus.Logger, err error, snapshotKey string, deleteSnapshotResult *report.Result) {
	if err != nil {
		msg := fmt.Sprintf("Failed to delete snapshot with key <%v>. Error <%v>", snapshotKey, err)
		logger.Error(msg)
		deleteSnapshotResult.Result = false
		deleteSnapshotResult.Cause = msg
	} else {
		deleteSnapshotResult.Result = true
		deleteSnapshotResult.Cause = ""
	}
}

func setDashBoardSnapshotResult(logger *logrus.Logger, err error, dashboardUID string, snapshotResult *report.Result) bool {

	if err != nil {
		logger.Errorf("Unable to create snapshot for dashboard <%v>. err <%v>", dashboardUID, err)
		snapshotResult.Result = false
		snapshotResult.Cause = err.Error()

		return true
	}

	snapshotResult.Result = true
	snapshotResult.Cause = ""
	return false
}

func setGetPanelIDsResult(logger *logrus.Logger, pErr error, extractPanelIDResult *report.Result, ids []grafana.PanelDescriptor, includeIDs, excludeIDs []int, dashJson json.RawMessage) bool {

	if pErr != nil {
		logger.Errorf("Unable to parse grafana dashboard API response. <%v>", pErr)
		extractPanelIDResult.Result = false
		extractPanelIDResult.Cause = pErr.Error()

		return true
	}

	if len(ids) == 0 {
		//No panels left to be recorded
		msg := fmt.Sprintf("No panels ids remaining after applying inclusion <%v> and exclusion <%v> lists to dashboard result <%v>", includeIDs, excludeIDs, dashJson)
		logger.Info(msg)
		extractPanelIDResult.Result = false
		extractPanelIDResult.Cause = msg

		return true
	}

	//Set postitive result
	extractPanelIDResult.Result = true
	extractPanelIDResult.Cause = ""

	return false
}

//setDashExistsResult - returns true if an error or negative result was detected
func setDashExistsResult(logger *logrus.Logger, dashErr error, dashboard common.GrafanaDashBoard, existsResult *report.Result, gdbExists bool) bool {
	if dashErr != nil {
		logger.Errorf("Internal error <%v> when checking if dashboard <%v> exists at host <%v> port <%v>", dashErr, dashboard.UID, dashboard.Host, dashboard.Port)
		existsResult.Result = false
		existsResult.Cause = dashErr.Error()
		return true
	}
	if !gdbExists {
		msg := fmt.Sprintf("Dashboard <%v> does not exist at host <%v> port <%v>", dashboard.UID, dashboard.Host, dashboard.Port)
		logger.Debug(msg)
		existsResult.Result = false
		existsResult.Cause = msg
		return true
	}

	//Set positive result
	existsResult.Result = true
	existsResult.Cause = ""

	return false
}
