package grafana

import (
	"encoding/json"
	"github.com/sajeevany/graph-snapper/internal/common"
	"github.com/sirupsen/logrus"
	"strings"
	"time"
)

type CreateSnapshotResponse struct {
	DeleteKey string `json:"deleteKey"`
	DeleteUrl string `json:"deleteUrl"`
	Key       string `json:"key"`
	Url       string `json:"url"`
}

//CreateSnapshot - Create snapshot for the specified dashboard. Returns
func CreateSnapshot(logger *logrus.Logger, host string, port int, user common.Basic, startTime, endTime time.Time, expiry time.Time, dashboard json.RawMessage) (CreateSnapshotResponse, error) {

	return CreateSnapshotResponse{}, nil
}

type GetSnapshotsResponse []GSnapshot

type GSnapshot struct {
	ID          string    `json:"snapshotId"`
	Name        string    `json:"name"`
	Key         string    `json:"key"`
	OrgID       int       `json:"orgId"`
	UserId      int       `json:"userId"`
	External    bool      `json:"external"`
	ExternalUrl string    `json:"externalUrl"`
	Expires     time.Time `json:"expires"`
	Created     time.Time `json:"created"`
	Updated     time.Time `json:"updated"`
}

func (gsr GetSnapshotsResponse) GetSnapshot(id string) (bool, GSnapshot){

	for _, snapshot := range gsr {
		if strings.ToLower(snapshot.ID) == strings.ToLower(id){
			return true, snapshot
		}
	}

	return false, GSnapshot{}
}

func GetSnapshots(logger *logrus.Logger, host string, port int, user common.Basic) (GetSnapshotsResponse, error) {

	return GetSnapshotsResponse{}, nil
}

