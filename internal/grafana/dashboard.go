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

	logger.Debugf("Starting dashboard exists check for uid <%v> <%v>:<%v>", uid, host, port)

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
		Timeout: 1500 * time.Millisecond,
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

	logger.Debugf("Completed dashboard exists check for uid <%v> <%v>:<%v>. Returning <%v>, <%+v>", uid, host, port, true, string(dash[dashboardKey]))
	return true, dash[dashboardKey], nil
}

//Simplified dashboard for marshalling desired panel
type dashboard struct {
	Panels []panel `json:"panels"`
}

type panel struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
}

//PanelDescriptor - panel description used to describe panel for upload operation
type PanelDescriptor struct {
	Title       string
	SnapshotURL string
	ID          int
}

//PanelDescriptors - Slice of PanelsDescriptors that is sortable by ID
type PanelDescriptors []PanelDescriptor

func (s PanelDescriptors) Len() int {
	return len(s)
}
func (s PanelDescriptors) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s PanelDescriptors) Less(i, j int) bool {
	return s[i].ID < s[j].ID
}

//DownloadedPanelDesc - panel descriptor that has been downloaded to a local directory
type DownloadedPanelDesc struct {
	PanelDescriptor
	DownloadDir string
}

func GetPanelsDescriptors(logger *logrus.Logger, msg json.RawMessage, includeIDs, excludeIDs []int) ([]PanelDescriptor, error) {

	logger.Debug("Starting GetPanelsDescriptors() and filtering based on inclusion/exclusion filters")
	defer logger.Debug("Completed GetPanelsDescriptors()")

	//Check if msg is non-zero
	if msg == nil || len(msg) == 0 {
		return make([]PanelDescriptor, 0), nil
	}

	//Parse raw json
	var dash dashboard
	if uErr := json.Unmarshal(msg, &dash); uErr != nil {
		return nil, uErr
	}

	//Get panels
	panels := make(map[int]PanelDescriptor)
	for _, v := range dash.Panels {
		panels[v.ID] = PanelDescriptor{
			Title: v.Title,
			ID:    v.ID,
		}
	}

	return filterPanels(panels, includeIDs, excludeIDs), nil
}

func filterPanels(panels map[int]PanelDescriptor, include, exclude []int) []PanelDescriptor {

	//Validate input
	if panels == nil {
		return nil
	}
	if len(panels) == 0 {
		return make([]PanelDescriptor, 0)
	}

	//filter panels. Inclusion list takes priority over exclusion list
	if len(include) > 0 {
		//restrictive include - schedule will only ever snapshot these panels. New panels will not be included
		pInc := make([]PanelDescriptor, 0)
		for _, v := range include {
			//If a value exists in the inclusion slice then add it to included panels slice
			if panelSnap, exists := panels[v]; exists {
				pInc = append(pInc, panelSnap)
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

func mapToSlice(panels map[int]PanelDescriptor) []PanelDescriptor {

	//Validate input
	if panels == nil {
		return nil
	}

	//Create slice to return
	slc := make([]PanelDescriptor, len(panels))

	i := 0
	for _, panel := range panels {
		slc[i] = panel
		i++
	}
	return slc
}
