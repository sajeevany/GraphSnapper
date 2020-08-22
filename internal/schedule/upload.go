package schedule

import (
	"fmt"
	"github.com/sajeevany/graph-snapper/internal/common"
	"github.com/sajeevany/graph-snapper/internal/confluence"
	"github.com/sajeevany/graph-snapper/internal/grafana"
	"github.com/sajeevany/graph-snapper/internal/report"
	"github.com/sirupsen/logrus"
)

func setupDatastoreAndUploadDashboardPanels(logger *logrus.Logger, rep map[string]report.ConfluenceStoreStages, stores DataStores, dashboardName, dashboardUID string, panelImages []grafana.DownloadedPanelDesc) error {

	logger.Debug("Starting image upload to datastores()")

	//Create and upload images to confluence page. This page will house all dashboard and panel sub pages
	rep = make(map[string]report.ConfluenceStoreStages, len(stores.ConfluencePages))
	for _, confluenceStoreDir := range stores.ConfluencePages {

		reprt := report.NewConfluenceStoreStages(len(panelImages))
		rep[confluenceStoreDir.PageID] = reprt

		//Check if the top level confluence page exists. If not create it. This page is expected to house all dashbooard image uploads
		//Nesting structure:
		//parentPage(determined by ID)
		//   -> panelPage(determined by {dashboard}_{panelID}_{panelName} (Every image will be uploaded to this page and added as an attachment
		exists, eErr := confluence.PageExistsByID(logger, confluenceStoreDir.PageID, confluenceStoreDir.Host, confluenceStoreDir.Port, common.Auth{
			Basic: confluenceStoreDir.User,
		})
		//User should specify where to store all the snapshot images. Otherwise they won't be able to find out where the data is stored
		if failed := setTopPageExistsCheckResult(logger, eErr, confluenceStoreDir, &reprt.TopPageExistsCheck, exists); failed {
			continue
		}

		//Use dashboard name to see if dashboard page exists since we don't want the user to have to have to premake everything.
		//This dashboard page is where each panel page will be nested. Use dashboardName_dashboardUID so that it's easy
		//for users to navigate the confluence page
		var dashboardPageID string
		expectedDashboardPageTitle := fmt.Sprintf("%s_%s", dashboardName, dashboardUID)
		dashboardPageID, exists, eErr = confluence.PageExistsByName(logger, expectedDashboardPageTitle, confluenceStoreDir.SpaceKey, confluenceStoreDir.PageID, confluenceStoreDir.User)
		if eErr != nil {
			//Set error and return. Request always refer to images for a particular dashboard. No point continuing
		}
		if !exists {
			var dCErr error
			if dashboardPageID, dCErr = confluence.CreateDashboardTitlePage(logger, expectedDashboardPageTitle, confluenceStoreDir.SpaceKey, confluenceStoreDir.PageID, confluenceStoreDir.User); dCErr != nil {
				//Unable to create the dashboard page the desired spot. Return with error
			}
		}

		//Traverse to panel storage page and update if it exists, otherwise create
		for _, panelImage := range panelImages {
			rep[confluenceStoreDir.PageID].SnapshotUploads[panelImage.Title] = report.SnapshotUpload{}
			if uploadErr := uploadPanelImage(logger, rep[confluenceStoreDir.PageID].SnapshotUploads[panelImage.Title], confluenceStoreDir.SpaceKey, dashboardPageID, confluenceStoreDir.User, panelImages); uploadErr != nil {
				continue
			}
		}
	}

	return nil
}

func uploadPanelImage(logger *logrus.Logger, upload report.SnapshotUpload, key string, id string, user common.Basic, images []grafana.DownloadedPanelDesc) error {
	//Upload image to page
	//Get current attachment ids
	//Generate one that does not exist. {dashboard}_{panelID}_{panelName}_{from_range}_{to_range}
	//Upload attachment and reserve ID
	//Get page contents.
	//Add image reference to top of page. If it is in an invalid format, then return with error
	//Edit page contents
	//Return
	return nil
}

//setTopPageExistsCheckResult - sets result of parent exists check. Returns true if an error occurred or the parent page doesn't exist. Use to continue/skip any further operations.
func setTopPageExistsCheckResult(logger *logrus.Logger, eErr error, parent ConfluencePage, pExistsResult *report.Result, exists bool) bool {

	if eErr != nil {
		msg := fmt.Sprintf("Error checking if confluence page with source id <%v> does not exist. <%v>", parent.PageID, eErr.Error())
		logger.Errorf(msg)
		pExistsResult.Result = false
		pExistsResult.Cause = msg

		return true
	}
	if !exists {
		pExistsResult.Result = false
		pExistsResult.Cause = "Page does not exist"
		logger.Debugf("Page <%v> does not exist", parent.PageID)

		return true
	}

	//Set passed result
	pExistsResult.Result = true
	pExistsResult.Cause = ""

	return false
}
