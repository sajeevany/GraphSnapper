package confluence

import (
	"github.com/sajeevany/graph-snapper/internal/common"
	"github.com/sirupsen/logrus"
)

func PageExistsByName(logger *logrus.Logger, title, spaceKey, parentPageID string, user common.Basic) (string, bool, error) {
	return "", false, nil
}

func CreateDashboardTitlePage(logger *logrus.Logger, title, spaceKey, parentPageID string, user common.Basic) (string, error) {
	return "", nil
}

func CreatePage(logger *logrus.Logger, title, spaceKey, parentPageID string, user common.Basic) (string, error) {
	return "", nil
}
