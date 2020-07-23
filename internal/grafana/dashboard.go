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
		Timeout: 15 * time.Millisecond,
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
	if readErr == nil {
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
