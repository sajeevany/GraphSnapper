package schedule

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/chromedp/chromedp"
	"github.com/sajeevany/graph-snapper/internal/common"
	"github.com/sajeevany/graph-snapper/internal/grafana"
	"github.com/sajeevany/graph-snapper/internal/report"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"time"
)

//Snapshots a grafana dashboard and stores the images in a temporary folder. Returns dashboard name, a map of image descriptors to image locations, and image storage directory.
func captureDashboardPanels(logger *logrus.Logger, dashboardStages *report.GrafanaDBSnapshotStages, dashboard common.GrafanaDashBoard, downloadDir string) (string, []grafana.DownloadedPanelDesc, error) {

	//check if specified dashboard exists. Get the dashboard information so it can be reused to create the snapshot
	gdbExists, dashJson, dashErr := grafana.DashboardExists(logger, dashboard.UID, dashboard.Host, dashboard.Port, dashboard.User.Auth.Basic)
	if failureOccurred := setDashExistsResult(logger, dashErr, dashboard, &dashboardStages.DashboardExistsCheck, gdbExists); failureOccurred {
		return "", nil, dashErr
	}

	//Get the dashboard title from the GET dashboard by uid request. No report stage since this stage should not fail if
	//the dashboard exists check doesn't fail
	dashboardName, uErr := grafana.ExtractDashboardTitleFromGetDBReq(dashJson)
	if uErr != nil || dashboardName == "" {
		logger.Warn("Unable to extract dashboard title from <%v>. Received error <%v> or empty dashboard name", string(dashJson), uErr)
		return "", nil, uErr
	}

	//Get panels to be screencaptured. Skip dashboard if no panels are to be captured
	panelDesc, pErr := grafana.GetPanelsDescriptors(logger, dashJson, dashboard.IncludePanelsIDs, dashboard.ExcludePanelsIDs)
	if setFailedResult := setGetPanelIDsResult(logger, pErr, &dashboardStages.ExtractPanelID, panelDesc, dashboard.IncludePanelsIDs, dashboard.ExcludePanelsIDs, dashJson); setFailedResult {
		return "", nil, dashErr
	}

	//Create snapshot, screen capture and save to temporary directory
	images, cdErr := createAndDownloadSnapshotPanels(logger, dashboardStages, dashboard.Host, dashboard.Port, dashboard.User.Auth.Basic, dashboard.UID, dashJson, panelDesc, downloadDir)
	if cdErr != nil {
		//c&d method is responsible for updating the dashboard report
		logger.Debugf("An error occurred while creating and downloading the dashboard snapshot <%v>", cdErr)
		return "", nil, cdErr
	}
	logger.Debugf("Downloaded images as <%+v>", images)

	//upload all attachments with name to page
	return dashboardName, images, nil
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

func setPanelSnapshotUrls(host string, port int, snapshotID string, panelIDs []grafana.PanelDescriptor) []grafana.PanelDescriptor {
	for idx, panel := range panelIDs {
		panelIDs[idx].SnapshotURL = fmt.Sprintf(SnapshotURLFmt, host, port, snapshotID, panel.ID)
	}
	return panelIDs
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
