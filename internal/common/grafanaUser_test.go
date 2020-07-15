package common

import (
	"testing"
)

func TestAddGrafanaReadUserV1_IsValid(t *testing.T) {
	tests := []struct {
		name string
		user GrafanaUserV1
		want bool
	}{
		{
			name: "test0 empty API key",
			user: GrafanaUserV1{
				Auth: Auth{
					BearerToken: BearerToken{},
					Basic:       Basic{},
				},
				Host:        "10.2.3.4",
				Port:        9000,
				Description: "blah",
			},
			want: false,
		},
		{
			name: "test1 empty hostkey",
			user: GrafanaUserV1{
				Auth: Auth{
					BearerToken: BearerToken{
						Token: "abcdefg",
					},
					Basic: Basic{
						Username: "asc",
						Password: "qwerty",
					},
				},
				Host:        "",
				Port:        9000,
				Description: "blah",
			},
			want: false,
		},
		{
			name: "test2 0 val port",
			user: GrafanaUserV1{
				Auth: Auth{
					BearerToken: BearerToken{
						Token: "abcdefg",
					},
					Basic: Basic{
						Username: "asc",
						Password: "qwerty",
					},
				},
				Host:        "10.2.3.4",
				Port:        0,
				Description: "blah",
			},
			want: false,
		},
		{
			name: "test3 negative port",
			user: GrafanaUserV1{
				Auth: Auth{
					BearerToken: BearerToken{
						Token: "abcdefg",
					},
					Basic: Basic{
						Username: "asc",
						Password: "qwerty",
					},
				},
				Host:        "10.2.3.4",
				Port:        -1,
				Description: "blah",
			},
			want: false,
		},
		{
			name: "test4 over max port value",
			user: GrafanaUserV1{
				Auth: Auth{
					BearerToken: BearerToken{
						Token: "abcdefg",
					},
					Basic: Basic{
						Username: "asc",
						Password: "qwerty",
					},
				},
				Host:        "10.2.3.4",
				Port:        999999999999,
				Description: "blah",
			},
			want: false,
		},
		{
			name: "test5 valid entry",
			user: GrafanaUserV1{
				Auth: Auth{
					BearerToken: BearerToken{
						Token: "abcdefg",
					},
					Basic: Basic{
						Username: "asc",
						Password: "qwerty",
					},
				},
				Host:        "10.2.3.4",
				Port:        8090,
				Description: "blah",
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.user.IsValid(); got != tt.want {
				t.Errorf("IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}
