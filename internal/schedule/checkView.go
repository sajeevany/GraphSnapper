package schedule

import "github.com/sajeevany/graph-snapper/internal/common"

type CheckScheduleV1 struct {
	GraphDashboards DashBoards
	DataStores      DataStores
}

type DashBoards struct {
	GrafanaDashboards GrafanaDashboards
}

type GrafanaDashboards struct {
	Dashboards map[string]GrafanaDashBoard
	User       common.GrafanaUserV1
}

type GrafanaDashBoard struct {
	Host   string
	Port   string
	UID    string
	Panels map[string]Panel //if empty include all panels, if non empty only do these panels
}

type Panel struct {
	ID    string
	Title string
}

type DataStores struct {
	ConfluencePages map[string]ConfluencePage
}

//ConfluencePage defines the location in which pages will be created
type ConfluencePage struct {
	SpaceKey     string
	ParentPageID string
	User         common.ConfluenceServerUserV1
}

type TestResponse struct {
	DatastoreUrls []string
}
