package grafana

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/sajeevany/graph-snapper/internal/common"
	"github.com/sirupsen/logrus"
	"github.com/valyala/fastjson"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	minimumExpirationTime = 15 * time.Minute
	createSnapshotTimeout = 150 * time.Millisecond

	PostSnapshotURL = "http://%s:%d/api/snapshots"
	GetSnapshotsURL = "http://%s:%d/api/dashboard/snapshots"
)

type CreateSnapshotRequest struct {
	//https://grafana.com/docs/grafana/latest/http_api/snapshot/
	Dashboard json.RawMessage `json:"dashboard"`
	Name      string          `json:"name,omitempty"`
	Expires   int             `json:expires,omitempty`
	External  bool            `json:external,omitempty`
	Key       string          `json:key,omitempty`
	DeleteKey string          `json:deleteKey,omitempty` //Unique key that ensures only snapshot creator can delete it
}

func (r CreateSnapshotRequest) GetFields() logrus.Fields {
	return logrus.Fields{
		"dashboard": r.Dashboard,
		"name":      r.Name,
		"expires":   r.Expires,
		"external":  r.External,
		"key":       r.Key,
		"deleteKey": r.DeleteKey,
	}
}

type CreateSnapshotResponse struct {
	DeleteKey string `json:"deleteKey"`
	DeleteUrl string `json:"deleteUrl"`
	Key       string `json:"key"`
	Url       string `json:"url"`
}

//CreateSnapshot - Create snapshot for the specified dashboard. Assumes specified dashboard exists on the target machine
func CreateSnapshot(logger *logrus.Logger, host string, port int, user common.Basic, startTime, endTime time.Time, expiry time.Time, dashboard json.RawMessage) (CreateSnapshotResponse, error) {

	logger.WithFields(user.GetFields()).Debugf("Starting create snapshot for host <%v> port <%v>", host, port)

	//create expiry. If net result in seconds is less than minimum required time, then default to the minimum
	expirationInSeconds := validateAndGetExpiration(expiry, time.Now(), minimumExpirationTime)

	//set time range in dashboard
	if startTime.After(endTime) {
		msg := fmt.Sprintf("Start time <%v> occurs after or is the same as end time <%v>", startTime, endTime)
		logger.Error(msg)
		return CreateSnapshotResponse{}, fmt.Errorf(msg)
	}
	dashWithTimeRange, setErr := setTimeRange(dashboard, startTime, endTime)
	if setErr != nil {
		logger.Errorf("Unable to update dashboard json message with start and end times. err <%v>", setErr)
		return CreateSnapshotResponse{}, setErr
	}

	//Create request
	csr := CreateSnapshotRequest{
		Dashboard: dashWithTimeRange,
		Expires:   expirationInSeconds,
		External:  false,
	}
	req, rErr := BuildSnapshotRequest(logger, host, port, user, csr)
	if rErr != nil {
		logger.Errorf("Error building snapshot request <%v>", rErr)
		return CreateSnapshotResponse{}, nil
	}

	//Create client and execute request
	client := http.Client{
		Timeout: createSnapshotTimeout,
	}
	resp, rpErr := client.Do(req)
	if rpErr != nil {
		logger.Errorf("Error when executing request <%v>", req)
		return CreateSnapshotResponse{}, rpErr
	}
	defer resp.Body.Close()

	//Process request
	if resp.StatusCode != http.StatusOK {
		respErr := fmt.Errorf("received non 200 response <%v> after executing request <%+v>", resp.StatusCode, req)
		logger.Error(resp)
		return CreateSnapshotResponse{}, respErr
	}

	data, bErr := ioutil.ReadAll(resp.Body)
	if bErr != nil {
		logger.Errorf("Error reading from response body <%v>. err <%v>", resp.Body, bErr)
		return CreateSnapshotResponse{}, bErr
	}
	if data == nil || len(data) == 0 {
		nodataErr := fmt.Errorf("no data recevied from create snapshot request <%v>", req)
		logger.Error(nodataErr)
		return CreateSnapshotResponse{}, nodataErr
	}
	var snapResp CreateSnapshotResponse
	if uErr := json.Unmarshal(data, &snapResp); uErr != nil {
		logger.Errorf("Error unmarshalling create snapshot response body <%v> for request <%v>. err <%v>", string(data), req, uErr)
		return CreateSnapshotResponse{}, uErr
	}

	return snapResp, nil
}

func setTimeRange(message json.RawMessage, startTime, endTime time.Time) ([]byte, error) {

	b := string(message)

	var p fastjson.Parser
	v, err := p.Parse(b)
	if err != nil {
		return nil, err
	}

	//set values
	setStartErr := setTime(startTime, p, v, "time.from")
	if setStartErr != nil {
		return nil, setStartErr
	}
	setEndErr := setTime(endTime, p, v, "time.to")
	if setEndErr != nil {
		return nil, setEndErr
	}

	//get value as json bytes
	return v.MarshalTo([]byte{}), nil
}

func setTime(tyme time.Time, p fastjson.Parser, v *fastjson.Value, key string) error {
	startUnixStr := strconv.FormatInt(tyme.Unix(), 10)
	startUnix, pErr := p.Parse(startUnixStr)
	if pErr != nil {
		return pErr
	}
	v.Set(key, startUnix)
	return nil
}

func BuildSnapshotRequest(logger *logrus.Logger, host string, port int, user common.Basic, csr CreateSnapshotRequest) (*http.Request, error) {

	//Marshal request body
	csrBytes, mErr := json.Marshal(csr)
	if mErr != nil {
		logger.WithFields(csr.GetFields()).Error("Unable to marshal CreateSnapshotRequest")
		return nil, mErr
	}

	requestUrl := fmt.Sprintf(PostSnapshotURL, host, port)
	req, err := http.NewRequest(http.MethodPost, requestUrl, bytes.NewReader(csrBytes))
	if err != nil {
		logger.Errorf("Unable to create get request to <%v>. err <%v>", requestUrl, err)
		return nil, err
	}
	req.SetBasicAuth(user.Username, user.Password)
	req.Header.Add("Content-Type", "application/json")

	return req, nil
}

//Get the expiration in seconds. Returns the minimum if the expiry is below the required threshold
func validateAndGetExpiration(expiry, currentTime time.Time, minimumExpiry time.Duration) int {
	expirationDiff := expiry.Sub(currentTime)
	if expirationDiff < minimumExpiry {
		//round to int as that's what grafana takes
		return int(minimumExpiry.Seconds())
	}
	return int(expirationDiff.Seconds())
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

//GetSnapshot - Gets snapshot based on ID. Returns truthy based on existence and the snapshot if it exists
func (gsr GetSnapshotsResponse) GetSnapshot(id string) (bool, GSnapshot) {

	for _, snapshot := range gsr {
		if strings.ToLower(snapshot.ID) == strings.ToLower(id) {
			return true, snapshot
		}
	}

	return false, GSnapshot{}
}

func GetSnapshots(logger *logrus.Logger, host string, port int, user common.Basic) (GetSnapshotsResponse, error) {

	logger.Debugf("Starting GetSnapshots() for host <%v> port <%v>", host, port)

	//Create request
	url := fmt.Sprintf(GetSnapshotsURL, host, port)
	req, rErr := http.NewRequest(http.MethodGet, url, nil)
	if rErr != nil {
		logger.Errorf("Error creating get snapshot request as <%v>. error <%v>", url, rErr)
		return GetSnapshotsResponse{}, rErr
	}
	req.SetBasicAuth(user.Username, user.Password)

	//Create client and fire request
	client := http.Client{
		Timeout: 150 * time.Millisecond,
	}
	resp, resErr := client.Do(req)
	if resErr != nil {
		logger.Errorf("Error sending get snapshot request as <%v>. error <%v>", url, resErr)
		return GetSnapshotsResponse{}, resErr
	}
	defer resp.Body.Close()

	//Parse body
	data, dErr := ioutil.ReadAll(resp.Body)
	if dErr != nil {
		logger.Errorf("Error reading get snapshot request response from <%v>. error <%v>", url, dErr)
		return GetSnapshotsResponse{}, dErr
	}

	//Marshal value and return
	var snapResp GetSnapshotsResponse
	if uErr := json.Unmarshal(data, &snapResp); uErr != nil {
		logger.Errorf("Error unmarshalling get snapshot request response from <%v>. error <%v>", url, uErr)
	}

	return snapResp, nil
}

func DeleteSnapshot(logger *logrus.Logger, host string, port int, user common.Basic, deleteKey string) error {
	return nil
}
