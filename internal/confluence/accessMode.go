package confluence

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"strings"
)

const AccessModeURL = "/rest/api/accessmode"

func HasWriteAccess(logger *logrus.Logger, host string, port int, user, password string) (bool, error) {

	logger.Debug("Starting valid login API key check")

	client := http.Client{}
	req, err := buildAccessModeRequest(logger, host, port, user, password)
	if err != nil {
		logger.Debugf("An error was found when creating http request to validate confluence user. <%v>", err)
		return false, err
	}

	//execute
	resp, rErr := client.Do(req)
	if rErr != nil {
		logger.Debugf("Error when calling request to <%v>. err <%v>", req.URL, err)
		return false, rErr
	}
	logger.Debug("accessMode API request executed. Checking status code.")
	defer resp.Body.Close()

	//Check response
	switch resp.StatusCode {
	case http.StatusOK:

		logger.Debug("confluence user access mode check returned 200")
		//Check if the resultant json string contains "write"
		hasWAccess, err2 := respHasWrite(resp)
		if err2 != nil {
			logger.Debugf("Error when reading response body <%v>, err <%v>", resp.Body, err2)
			return hasWAccess, err2
		}

		return hasWAccess, nil
	default:
		logger.Debugf("Unexpected response status code <%v>", resp.StatusCode)
		return false, nil
	}

}

func buildAccessModeRequest(logger *logrus.Logger, host string, port int, user, password string) (*http.Request, error) {

	//Build request url
	reqURL := fmt.Sprintf("http://%v:%v%v", host, port, AccessModeURL)

	//Create request
	req, err := http.NewRequest(http.MethodGet, reqURL, nil)
	if err != nil {
		logger.Debugf("An error was found when creating http request to validate confluence user. <%v>", err)
		return nil, err
	}

	//add headers
	req.SetBasicAuth(user, password)
	req.Header.Set("Content-Type", "application/json")

	return req, nil
}

//Checks if response has the word write in it
func respHasWrite(resp *http.Response) (bool, error) {

	bodyBytes, bErr := ioutil.ReadAll(resp.Body)
	if bErr != nil {
		return false, bErr
	}
	body := strings.ToLower(string(bodyBytes))
	hasWAccess := strings.Contains(body, "write")

	return hasWAccess, nil
}
