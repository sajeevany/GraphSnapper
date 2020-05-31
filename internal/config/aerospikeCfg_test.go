package config

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"reflect"
	"testing"
)

//Tests aerospikeCfg.addInvalidArg
func TestAerospikePortfolioConfig_AddInvalidArg(t *testing.T) {

	type setup struct{
		jsonPath string
		asConf AerospikeCfg
	}

	type expectedResult struct{
		ok bool
		invalidArgs map[string]string
	}

	// Testing scenarios
	var tests = []struct {
		expectedResult expectedResult
		setup setup
	}{
		{
			//All attributes are invalid
			expectedResult: expectedResult{

				ok:          false,
				invalidArgs: map[string]string{
					"conf.aerospike.host" : fmt.Sprintf("<%v> field is using an invalid value <%v>", "host", ""),
					"conf.aerospike.password" : fmt.Sprintf("<%v> field is using an invalid value <%v>", "password", ""),
					"conf.aerospike.port" : fmt.Sprintf("<%v> field is using an invalid value <%v>", "port", "0"),
					"conf.aerospike.graphNamespace" : fmt.Sprintf("<%v> field is using an invalid value <%v>", "graphNamespace", ""),
				},
			},
			setup: setup{
				jsonPath: "conf.aerospike",
				asConf: AerospikeCfg{
					Host:           "", // Cannot be empty
					Port:           0,  // Cannot be zero
					Password:       "", // Cannot be empty
					GraphNamespace: "", // Cannot be empty
				},
			},
		},
	}

	//Nil logger to be used in testing
	logger := logrus.New()


	// Execute test
	for _, scenario := range tests{

		//scenario attributes
		sSetup := scenario.setup
		sExpect := scenario.expectedResult


		//run valid check
		ok, invalidArgs := sSetup.asConf.IsValid(logger, sSetup.jsonPath)

		//compare with expected result
		if ok != sExpect.ok{
			t.Errorf("Expected <%v> for ok but was <%v>", sExpect.ok, ok)
		}

		if !reflect.DeepEqual(sExpect.invalidArgs, invalidArgs){
			actual, _ := json.MarshalIndent(invalidArgs, "", "\t")
			expected , _ := json.MarshalIndent(sExpect.invalidArgs, "", "\t")
			t.Errorf("Actual result <%v> does not match expected value <%v>", string(actual), string(expected))
		}

	}

}