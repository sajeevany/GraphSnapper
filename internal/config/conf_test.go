package config

import (
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
		invalidArgs []string
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
				invalidArgs: []string{
					"conf.aerospike.Host",
					"conf.aerospike.Port",
					"conf.aerospike.ConnectionRetries",
					"conf.aerospike.AccountNamespace",
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
				invalidArgs: []string{
					"conf.aerospike.AccountNamespace.Namespace",
				},
			},
			setup: setup{
				jsonPath: "conf.aerospike.AccountNamespace",
				asConf: AerospikeNamespace{
					Namespace: "", // Cannot be empty
					SetName:   "blah",
				},
			},
		},
		{
			testName: "TestAerospikePortfolioConfig_AddInvalidArg_2: logging level is inccorrect",
			expectedResult: expectedResult{
				ok: false,
				invalidArgs: []string{
					"conf.logging.Level",
				},
			},
			setup: setup{
				jsonPath: "conf.logging",
				asConf:   Logging{Level: "car"},
			},
		},
	}

	// Execute testName
	for _, scenario := range tests {

		//scenario attributes
		sSetup := scenario.setup
		sExpect := scenario.expectedResult

		//run valid check
		invalidArgs := make(map[string]string)
		ok := sSetup.asConf.IsValid(sSetup.jsonPath, invalidArgs)

		//compare with expected result
		if ok != sExpect.ok {
			t.Errorf("Expected asConf.IsValid for to be  <%v> but was  <%v> ", sExpect.ok, ok)
			t.Fail()
		}

		//Compare map size to invalid arg keys
		if len(invalidArgs) != len(sExpect.invalidArgs) {
			t.Errorf("Expected number of invalid arguments <%v> to match expected number of invalid args <%v>. Expected <%v> Actual <%v>", len(sExpect.invalidArgs), len(invalidArgs), sExpect.invalidArgs, invalidArgs)
			t.Fail()
		}

		//Iterate over expected invalid keys and assert that they exist in map
		for _, v := range sExpect.invalidArgs {
			if _, exists := invalidArgs[v]; !exists {
				t.Errorf("Expected arg <%v> in <%v> but was not found", v, invalidArgs)
				t.FailNow()
			}
		}

		if t.Failed() {
			t.Logf("Failed test description: %v", scenario.testName)
		}
	}

}
