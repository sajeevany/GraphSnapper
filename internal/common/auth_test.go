package common

import (
	"encoding/base64"
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
	"reflect"
	"testing"
)

func TestSetAuthHeader(t *testing.T) {
	type args struct {
		auth             Auth
		expectedAuthType string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Test 0: Only basic auth is set. Expect basic header",
			args: args{
				auth: Auth{
					BearerToken: BearerToken{},
					Basic: Basic{
						Username: "userName",
						Password: "password",
					},
				},
				expectedAuthType: BasicAuthType,
			},
		},
		{
			name: "Test 1: Bearer token is set. Expect bearer header",
			args: args{
				auth: Auth{
					BearerToken: BearerToken{
						Token: "tolkien",
					},
					Basic: Basic{},
				},
				expectedAuthType: BearerTokenAuthType,
			},
		},
		{
			name: "Test 2: Both are set. Expect bearer header",
			args: args{
				auth: Auth{
					BearerToken: BearerToken{
						Token: "tolkien",
					},
					Basic: Basic{
						Username: "userName",
						Password: "password",
					},
				},
				expectedAuthType: BearerTokenAuthType,
			},
		},
		{
			name: "Test 3: Neither are set. Expect no header",
			args: args{
				auth: Auth{
					BearerToken: BearerToken{},
					Basic:       Basic{},
				},
				expectedAuthType: "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "URL", nil)
			SetAuthHeader(logrus.StandardLogger(), tt.args.auth, req)

			switch tt.args.expectedAuthType {
			case BasicAuthType:
				//Grab base64 encoded result
				basicVal := req.Header.Values("Authorization")

				//Convert setup arguments to expected convert format. str -> base64 + prefix
				up := []byte(fmt.Sprintf("%s:%s", tt.args.auth.Basic.Username, tt.args.auth.Basic.Password))
				upBase64 := fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString(up))
				want := []string{upBase64}

				//Test
				if !reflect.DeepEqual(basicVal, want) {
					t.Errorf("beaterValue = %v, want %v", basicVal, want)
				}
			case BearerTokenAuthType:

				//Grab base64 encoded result
				bearerVal := req.Header.Values("Authorization")

				//Convert setup arguments to expected convert format. str -> base64 + prefix
				want := []string{"Bearer " + tt.args.auth.BearerToken.Token}

				//Test
				if !reflect.DeepEqual(bearerVal, want) {
					t.Errorf("beaterValue = %v, want %v", bearerVal, want)
				}
			default:

				authVal := req.Header.Values("Authorization")
				var want []string
				if !reflect.DeepEqual(authVal, want) {
					t.Errorf("beaterValue = %v, want <%v>", authVal, want)
				}
			}

		})
	}
}
