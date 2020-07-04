package grafana

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
)

const loginPingURL = "/api/login/ping"

func IsValidLogin(logger *logrus.Logger, apiKey, host string, port int) (bool, error) {

	logger.Debug("Starting valid login API key check")

	client := http.Client{}
	reqURL := buildLoginRequestURL(host, port)
	req, err := http.NewRequest(http.MethodGet, reqURL, nil)
	if err != nil {
		logger.Debugf("An error was found when creating http request to validate grafana user. <%v>", err)
		return false, err
	}

	//add headers
	setAuthHeader(req, apiKey)
	logger.Debug("auth headers set")

	//execute
	resp, rErr := client.Do(req)
	if rErr != nil {
		logger.Debugf("Error when calling request to <%v>. err <%v>", reqURL, err)
		return false, rErr
	}
	logger.Debug("login API request executed. Checking status code.")
	defer resp.Body.Close()

	//Check response
	switch resp.StatusCode {
	case http.StatusOK:
		logger.Debug("grafana user check returned 200")
		return true, nil
	case http.StatusUnauthorized:
		logger.Debugf("Unauthorized (401) response body <%v>", resp.Body)
		return false, nil
	default:
		logger.Debug("Unexpected response status code <%v>", resp.StatusCode)
		return false, nil
	}
}

func buildLoginRequestURL(host string, port int) string {
	return fmt.Sprintf("http://%v:%v%v", host, port, loginPingURL)
}

func setAuthHeader(req *http.Request, apikey string) {
	bearer := fmt.Sprintf("Bearer %v", apikey)
	req.Header.Add("Authorization", bearer)
}
