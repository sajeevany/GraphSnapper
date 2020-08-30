package confluence

import (
	"fmt"
	"github.com/sajeevany/graph-snapper/internal/common"
	"github.com/sirupsen/logrus"
)

//PageExistsByNameUnderID - check if a page with a particular title exists under a given space key with an input ancestor.
//Returns exists result, page id if page exists, contents as a raw json, and an error if one occurs
func PageExistsByNameUnderID(logger *logrus.Logger, title, spaceKey, pageAncestor string, user common.Basic) (bool, string, string, error) {

	logger.Debug("Starting PageExistsByNameUnderID check for <%s> page name in space <%s>")

	if err := validateInputPageExistsArgs(title, spaceKey, pageAncestor, user); err != nil {
		logger.Debugf("Invalid arguments when checking if PageExistsByNameUnderID. Input values may be empty or invalid. error <%v>", err)
		return false, "", "", err
	}

	return false, "", "", nil
}

func validateInputPageExistsArgs(title string, key string, ancestor string, user common.Basic) error {

	if title != "" && key != "" && ancestor != "" && user.IsValid() {
		return nil
	}

	return fmt.Errorf("")
}

func CreateDashboardTitlePage(logger *logrus.Logger, title, spaceKey, parentPageID string, user common.Basic) (string, error) {
	return "", nil
}

func CreatePage(logger *logrus.Logger, title, spaceKey, parentPageID string, user common.Basic) (string, error) {
	return "", nil
}
