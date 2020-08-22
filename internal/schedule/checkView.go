package schedule

import (
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
