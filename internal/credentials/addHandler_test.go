package credentials

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
	"github.com/sajeevany/graph-snapper/internal/account"
	"github.com/sajeevany/graph-snapper/internal/common"
	"github.com/sajeevany/graph-snapper/internal/db/aerospike"
	"github.com/sajeevany/graph-snapper/internal/db/aerospike/record"
	"github.com/sajeevany/graph-snapper/internal/logging"
	"github.com/sajeevany/graph-snapper/internal/test"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"
)

//PutCredentialsIntegrationTest
func TestPutCredentialsV1Integration(t *testing.T) {

	//Skip test if user wants to only run regression tests
	if testing.Short() {
		t.Skip()
	}

	//Setup common requirements. In this case it's a specific aerospike image.
	ctx := context.Background()
	aeroContainer, aeroClient := test.StartAerospikeTestContainer(t, ctx)
	defer aeroContainer.Terminate(ctx)

	type expected struct {
		returnCode int
		creds      record.CredentialsView1
	}

	//Scenarios
	tests := []struct {
		name      string
		setup     func(logger *logrus.Logger, client *aerospike.ASClient, accountKey string)
		cleanup   func(client *aerospike.ASClient)
		accountID string
		request   SetCredentialsV1
		expected  expected
	}{
		{
			name: "test0 PutCredentialsV1 happy path",
			setup: func(logger *logrus.Logger, client *aerospike.ASClient, accountKey string) {
				recReq := record.AccountViewV1{
					Email: "testUser@graphSnapper.com",
					Alias: "Admin config account",
				}
				expectedAct := record.AccountV1{
					Email: "testUser@graphSnapper.com",
					Alias: "Admin config account",
				}
				//Create account
				rec, err := account.CreateAccount(logger, client, accountKey, recReq)
				if err != nil {
					t.Errorf("SETUP FAILURE: An error occurred when creating a new account record <%#v>, err <%v>", recReq, err)
				}
				if !reflect.DeepEqual(rec.Account, expectedAct) {
					t.Errorf("SETUP FAILURE: Create account operation did not create an account as expected. Expected <%+v>\n Actual <%+v>\n", recReq, expectedAct)
				}
			},
			cleanup: func(asClient *aerospike.ASClient) {
				ns := asClient.AccountNamespace
				tyme := time.Now()
				if err := asClient.Client.Truncate(nil, ns.Namespace, ns.SetName, &tyme); err != nil {
					t.Errorf("CLEANUP FAILURE: Unable to truncate test aerospike container namespace <%v>, err <%v>", ns, err)
				}
			},
			accountID: "abc",
			request: SetCredentialsV1{
				GrafanaAPIUsers: map[string]common.GrafanaUserV1{
					"gu_0": {
						Auth: common.Auth{
							BearerToken: common.BearerToken{
								Token: "gu0APIToken",
							},
							Basic: common.Basic{},
						},
						Host:        "test0.grafanahost.com",
						Port:        8565,
						Description: "test0 grafana auth",
					},
				},
				ConfluenceServerUsers: map[string]common.ConfluenceServerUserV1{
					"csu_0": {
						Host:        "test0.host.com",
						Port:        9220,
						Description: "test0 confluence",
						Auth: common.Auth{
							Basic: common.Basic{
								Username: "confluenceUsername",
								Password: "confluencePassword",
							},
						},
					},
				},
			},
			expected: expected{
				returnCode: 200,
				creds: record.CredentialsView1{
					GrafanaAPIUsers: map[string] record.GrafanaAPIUser{
						"gu_0": {
							Auth: common.Auth{
								BearerToken: common.BearerToken{
									Token: logging.RedactNonEmpty("gu0APIToken"),
								},
								Basic: common.Basic{},
							},
							Host:        "test0.grafanahost.com",
							Port:        8565,
							Description: "test0 grafana auth",
						},
					},
					ConfluenceServerUsers: map[string]record.ConfluenceServerUser{
						"csu_0": {
							Host:        "test0.host.com",
							Port:        9220,
							Description: "test0 confluence",
							Auth: common.Auth{
								Basic: common.Basic{
									Username: logging.RedactNonEmpty("confluenceUsername"),
									Password: logging.RedactNonEmpty("confluencePassword"),
								},
							},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//setup and queue cleanup
			accountID := "abcde"
			logger := logrus.New()
			tt.setup(logger, aeroClient, accountID)
			defer tt.cleanup(aeroClient)

			//Build request
			j, mErr := jsoniter.Marshal(tt.request)
			if mErr != nil {
				t.Errorf("Error marshalling request <%+v>", tt.request)
			}
			req, rErr := http.NewRequest(http.MethodPut, fmt.Sprintf("/api/v1/account/%s/credentials", tt.accountID), bytes.NewBuffer(j))
			if rErr != nil {
				t.Errorf("Error creating new request")
			}
			req.Header.Add("Content-Type", "application/json")

			//Setup gin engine to receive requests
			w := httptest.NewRecorder()
			gin.SetMode(gin.TestMode)
			_, r := gin.CreateTestContext(w)
			r.PUT("/api/v1/account/:id/credentials", PutCredentialsV1(logger, aeroClient))

			//Run Test
			r.ServeHTTP(w, req)

			//Validate
			if w.Code != tt.expected.returnCode {
				t.Errorf("Incorrect return code. Expected <%v> got <%v>", w.Code, tt.expected.returnCode)
			}
			data, bErr := ioutil.ReadAll(w.Body)
			if bErr != nil || data == nil {
				t.Errorf("Unable to read from http response <%v>", bErr)
			}
			var creds record.RecordViewV1
			if uErr := json.Unmarshal(data, &creds); uErr != nil {
				t.Errorf("Unable to unmarshal response err <%v>", uErr)
			}
			if !reflect.DeepEqual(tt.expected.creds, creds.Credentials) {
				t.Errorf("AddedCredentialsResponse does not match expected response. Expected <%#v>\n Actual <%#v>", tt.expected.creds, spew.Sdump(creds.Credentials))

				t.Logf("Expected: %v",  spew.Sdump(tt.expected.creds))
				t.Logf("Actual %v",spew.Sdump(creds.Credentials))

			}
		})
	}
}
