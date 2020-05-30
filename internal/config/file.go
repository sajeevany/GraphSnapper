package config

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
)

type Conf struct {
	Version string `json:"version"`
	Port    int    `json:"port"`
}

func (c Conf) getFields() logrus.Fields {
	return logrus.Fields{"version": c.Version, "port": c.Port}
}

//Read - reads config file referenced by conf
func Read(conf string, logger *logrus.Logger) (Conf, error) {

	logger.Debugf("Checking if file <%v> exists", conf)

	if _, err := os.Stat(conf); err == nil {
		//file exists. Go forth and conquer

		//Read file contents
		data, err := ioutil.ReadFile(conf)
		if err != nil{
			defaultConf := getDefaultConf()
			logger.WithFields(defaultConf.getFields()).Errorf("Error reading configuration file <%v>. Returning default config <%v>. Encountered error <%v>", conf, defaultConf, err)
			return defaultConf, err
		}

		//Unmarshal data as json
		var cStruct Conf
		if convErr := json.Unmarshal(data, &cStruct); convErr != nil{
			defaultConf := getDefaultConf()
			logger.WithFields(defaultConf.getFields()).Errorf("Error unmarshalling configuration file <%v>. Returning defaults <%v>. Encountered error <%v>.", conf, defaultConf, convErr)
			return defaultConf, convErr
		}

		return cStruct, nil

	} else if os.IsNotExist(err) {
		//file doesn't exist
		defaultConf := getDefaultConf()
		logger.WithFields(defaultConf.getFields()).Errorf("Configuration file <%v> does not exist. Using defaults.  Encountered error <%v>", conf, err)
		return defaultConf, err
	} else{
		//https://stackoverflow.com/questions/12518876/how-to-check-if-a-file-exists-in-go
		defaultConf := getDefaultConf()
		logger.WithFields(defaultConf.getFields()).Errorf("Error while evaluating if config file <%v> exists.", conf)
		return getDefaultConf(), err
	}
}

func getDefaultConf() Conf {
	return Conf{
		Version: "Default",
		Port:    8080,
	}
}