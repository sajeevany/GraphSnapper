package config

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
)

//Read - reads config file referenced by conf
func Read(conf string, logger *logrus.Logger) (*Conf, error) {

	logger.Debugf("Checking if file <%v> exists", conf)

	if _, err := os.Stat(conf); err == nil {
		//file exists. Go forth and conquer

		//Read file contents
		data, err := ioutil.ReadFile(conf)
		if err != nil {
			logger.Errorf("Error reading configuration file <%v>. Encountered error <%v>", conf, err)
			return nil, err
		}

		//Unmarshal data as json
		var cStruct Conf
		if convErr := json.Unmarshal(data, &cStruct); convErr != nil {
			logger.Errorf("Error unmarshalling configuration file <%v>. Encountered error <%v>.", conf, convErr)
			return nil, convErr
		}

		return &cStruct, nil

	} else if os.IsNotExist(err) {
		//file doesn't exist
		logger.Errorf("Configuration file <%v> does not exist. Encountered error <%v>", conf, err)
		return nil, err
	} else {
		logger.Errorf("Error <%v> while evaluating if config file <%v> exists.", err, conf)
		return nil, err
	}
}
