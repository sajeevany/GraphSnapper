package grafana

import (
	"encoding/json"
	"fmt"
	"github.com/sajeevany/graph-snapper/internal/common"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	GetDashboardURL = "http://%s:%d/api/dashboards/uid/%s"

	//key to get dashboard json message from GET UID response
	dashboardKey = "dashboard"
)

//DashboardExists - Checks if a dashboard with the specified UID exists at the specified location. Returns a bool result,
//and the dashboard json description (can be used in a snapshot call)
func DashboardExists(logger *logrus.Logger, uid, host string, port int, user common.Basic) (bool, json.RawMessage, error) {

	//Create request
	requestUrl := fmt.Sprintf(GetDashboardURL, host, port, uid)
	req, err := http.NewRequest(http.MethodGet, requestUrl, nil)
	if err != nil {
		logger.Errorf("Unable to create get request to <%v>. err <%v>", requestUrl, err)
		return false, nil, err
	}
	req.SetBasicAuth(user.Username, user.Password)

	//create client with timeout
	client := &http.Client{
		Timeout: 150 * time.Millisecond,
	}

	//Send request
	resp, err := client.Do(req)
	if err != nil {
		logger.Errorf("Error sending and receiving request from <%v>. err <%v>", requestUrl, err)
		return false, nil, err
	}
	defer resp.Body.Close()

	//Check return code
	if resp.StatusCode != http.StatusOK {
		logger.Debugf("Grafana dashboard check returned non 200 return code <%v>", resp.StatusCode)
		return false, nil, nil
	}

	//Get body
	body, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		logger.Errorf("Unable to read response body from request <%v>. err <%v>", requestUrl, readErr)
		return false, nil, err
	}

	//Marshal result so the dashboard value can be extracted
	var dash map[string]json.RawMessage
	if uErr := json.Unmarshal(body, &dash); uErr != nil {
		logger.Errorf("Unable to unmarshal response body <%v> into map[string]json.RawMessage. err <%v>", string(body), uErr)
		return false, nil, uErr
	}

	return true, dash[dashboardKey], nil
}

//Simplified dashboard for marshalling desired pane snapshotId
type dashboard struct {
	Panels []panel `json:"panels"`
}

type panel struct {
	ID int `json:"id"`
}

func GetPanelsIDs(msg json.RawMessage, includeIDs, excludeIDs []int) ([]int, error) {

	//Check if msg is non-zero
	if msg == nil || len(msg) == 0 {
		return []int{}, nil
	}

	//Parse raw json
	var dash dashboard
	if uErr := json.Unmarshal(msg, &dash); uErr != nil {
		return nil, uErr
	}

	//Get panels
	panels := make(map[int]struct{})
	for _, v := range dash.Panels {
		panels[v.ID] = struct{}{}
	}

	return filterPanels(panels, includeIDs, excludeIDs), nil
}

func filterPanels(panels map[int]struct{}, include, exclude []int) []int {

	//Validate input
	if panels == nil {
		return nil
	}
	if len(panels) == 0 {
		return []int{}
	}

	//filter panels. Inclusion list takes priority over exclusion list
	if len(include) > 0 {
		//restrictive include - schedule will only ever snapshot these panels. New panels will not be included
		var pInc []int
		for _, v := range include {
			//If a value exists in the inclusion slice then add it to included panels slice
			if _, exists := panels[v]; exists {
				pInc = append(pInc, v)
			}
		}
		return pInc
	} else if len(exclude) > 0 {
		//restrictive exclude - schedule will only exclude these panels. New panels will be automatically included
		for _, v := range exclude {
			//Delete an ID from the map of panels if it exists in the exclusion slice
			if _, exists := panels[v]; exists {
				delete(panels, v)
			}
		}
		return mapToSlice(panels)
	}

	//If no inclusion/exclusion slices are provided then every panel is eligible
	return mapToSlice(panels)
}

func mapToSlice(panels map[int]struct{}) []int {

	//Validate input
	if panels == nil {
		return nil
	}

	//Create slice to return
	slc := make([]int, len(panels))

	ctr := 0
	for key := range panels {
		slc[ctr] = key
		ctr++
	}
	return slc
}
