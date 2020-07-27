package confluence

import (
	"github.com/sajeevany/graph-snapper/internal/common"
	"github.com/sajeevany/graph-snapper/internal/grafana"
	"github.com/sirupsen/logrus"
)

func CreatePage(logger *logrus.Logger, title, spaceKey, parentPageID string, user common.Basic, panels map[grafana.PanelDescriptor]string) (string, error) {
	return "", nil
}
