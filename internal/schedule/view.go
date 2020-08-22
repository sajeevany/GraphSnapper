package schedule

import (
	"github.com/sajeevany/graph-snapper/internal/common"
	"github.com/sirupsen/logrus"
)

type ScheduleV1 struct {
	GraphDashboards DashBoards
	DataStores      DataStores
	Cadence         RecordDashboardCadence
}

type DashBoards struct {
	Grafana map[string]common.GrafanaDashBoard
}

func (b DashBoards) GetFields() logrus.Fields {

	grafanadb := make(logrus.Fields, len(b.Grafana))
	for key, db := range b.Grafana {
		grafanadb[key] = db.GetFields()
	}

	return logrus.Fields{
		"Grafana": grafanadb,
	}
}

type DataStores struct {
	ConfluencePages map[string]ConfluencePage
}

func (s DataStores) GetFields() logrus.Fields {
	pConf := make(logrus.Fields, len(s.ConfluencePages))
	for key, cp := range s.ConfluencePages {
		pConf[key] = cp.GetFields()
	}

	return logrus.Fields{
		"ConfluencePage": pConf,
	}
}

//ConfluencePage defines the confluence page location
type ConfluencePage struct {
	SpaceKey string
	PageID   string
	Host     string
	Port     int
	User     common.Basic
}

func (p ConfluencePage) GetFields() logrus.Fields {
	return logrus.Fields{
		"SpaceKey": p.SpaceKey,
		"PageID":   p.PageID,
		"User":     p.User.GetFields(),
	}
}

type RecordDashboardCadence struct {
	TriggerCron string
	Ranges      DashBoardTimeRanges
}
