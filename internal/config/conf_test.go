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

	type setup struct {
		jsonPath string
		asConf   ValidatableConf
	}

	type expectedResult struct {
		ok          bool
		invalidArgs map[string]string
	}

	// Testing scenarios
	var tests = []struct {
		testName       string
		expectedResult expectedResult
		setup          setup
	}{
		{
			testName: "TestAerospikePortfolioConfig_AddInvalidArg_0: All attributes are invalid",
			//All attributes are invalid
			expectedResult: expectedResult{

				ok: false,
				invalidArgs: map[string]string{
					"conf.aerospike.Host":             fmt.Sprintf("<%v> field is using an invalid value <%v>", "Host", ""),
					"conf.aerospike.Port":             fmt.Sprintf("<%v> field is using an invalid value <%v>", "Port", "0"),
					"conf.aerospike.AccountNamespace": fmt.Sprintf("<%v> field is using an invalid value <%v>", "AccountNamespace", ""),
				},
			},
			setup: setup{
				jsonPath: "conf.aerospike",
				asConf: AerospikeCfg{
					Host:             "",                   // Cannot be empty
					Port:             0,                    // Cannot be zero
					AccountNamespace: AerospikeNamespace{}, // Cannot be empty or zero AerospikeNamespace
				},
			},
		},
		{
			testName: "TestAerospikePortfolioConfig_AddInvalidArg_1: conf.Account Namespace has missing namespace",
			//AccountNamespace is missing namespace name
			expectedResult: expectedResult{

				ok: false,
				invalidArgs: map[string]string{
					"conf.aerospike.AccountNamespace.namespace": fmt.Sprintf("<%v> field is using an invalid value <%v>", "namespace", ""),
				},
			},
			setup: setup{
				jsonPath: "conf.aerospike",
				asConf: AerospikeCfg{
					Host: "abc", // Cannot be empty
					Port: 8080,  // Cannot be zero
					AccountNamespace: AerospikeNamespace{
						Namespace: "", // Cannot be empty
						SetName:   "blah",
					},
				},
			},
		},
		{
			testName: "TestAerospikePortfolioConfig_AddInvalidArg_2: logging level is inccorrect",
			expectedResult: expectedResult{
				ok: false,
				invalidArgs: map[string]string{
					"conf.logging.Level": fmt.Sprintf("<%v> field is using an invalid value <%v>", "Level", "car"),
				},
			},
			setup: setup{
				jsonPath: "conf.logging",
				asConf:   Logging{Level: "car"},
			},
		},
	}

	//Nil logger to be used in testing
	logger := logrus.New()

	// Execute testName
	for _, scenario := range tests {

		//scenario attributes
		sSetup := scenario.setup
		sExpect := scenario.expectedResult

		//run valid check
		invalidArgs := make(map[string]string)
		ok := sSetup.asConf.IsValid(logger, sSetup.jsonPath, invalidArgs)

		//compare with expected result
		if ok != sExpect.ok {
			t.Errorf("Expected asConf.IsValid for to be  <%v> but was  <%v> ", sExpect.ok, ok)
			t.Fail()
		}

		if !reflect.DeepEqual(sExpect.invalidArgs, invalidArgs) {
			actual, _ := json.MarshalIndent(invalidArgs, "", "\t")
			expected, _ := json.MarshalIndent(sExpect.invalidArgs, "", "\t")
			t.Errorf("Actual result <%v> does not match expected value <%v>", string(actual), string(expected))
		}

		if t.Failed() {
			t.Logf("Failed test description: %v", scenario.testName)
		}
	}

}
