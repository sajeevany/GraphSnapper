package schedule

import (
	"github.com/sajeevany/graph-snapper/internal/common"
	"github.com/sirupsen/logrus"
)

type CheckScheduleV1 struct {
	GraphDashboards DashBoards
	DataStores      DataStores
}

func (v CheckScheduleV1) IsValid() (bool, error) {

	//TODO implement
	return true, nil
}

func (v CheckScheduleV1) GetFields() logrus.Fields {
	return logrus.Fields{
		"DashBoards": v.GraphDashboards.GetFields(),
		"DataStores": v.DataStores.GetFields(),
	}
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
	ConfluencePages map[string]ParentConfluencePage
}

func (s DataStores) GetFields() logrus.Fields {
	pConf := make(logrus.Fields, len(s.ConfluencePages))
	for key, cp := range s.ConfluencePages {
		pConf[key] = cp.GetFields()
	}

	return logrus.Fields{
		"ParentConfluencePage": pConf,
	}
}

//ParentConfluencePage defines the parent location in which pages will be created as sub-pages
type ParentConfluencePage struct {
	SpaceKey     string
	ParentPageID string
	User         common.Basic
}

func (p ParentConfluencePage) GetFields() logrus.Fields {
	return logrus.Fields{
		"SpaceKey":     p.SpaceKey,
		"ParentPageID": p.ParentPageID,
		"User":         p.User.GetFields(),
	}
}
