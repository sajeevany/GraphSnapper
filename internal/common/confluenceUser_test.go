package common

import (
	"testing"
)

func TestAConfluenceServerUserV1_IsValid(t *testing.T) {
	tests := []struct {
		name string
		user ConfluenceServerUserV1
		want bool
	}{
		{
			name: "test0 empty username key",
			user: ConfluenceServerUserV1{
				Authentication: Auth{
					BearerToken: BearerToken{},
					Basic:       Basic{
						Username: "",
						Password: "qwerty",
					},
				},
				Host:     "10.2.3.4",
				Port:     9000,
			},
			want: false,
		},
		{
			name: "test1 empty password",
			user: ConfluenceServerUserV1{
				Authentication: Auth{
					BearerToken: BearerToken{},
					Basic:       Basic{
						Username: "user",
						Password: "",
					},
				},
				Host:     "10.2.3.4",
				Port:     9000,
			},
			want: false,
		},
		{
			name: "test2 0 val port",
			user: ConfluenceServerUserV1{
				Authentication: Auth{
					BearerToken: BearerToken{},
					Basic:       Basic{
						Username: "asc",
						Password: "qwerty",
					},
				},
				Host:     "10.2.3.4",
				Port:     0,
			},
			want: false,
		},
		{
			name: "test3 negative port",
			user: ConfluenceServerUserV1{
				Authentication: Auth{
					BearerToken: BearerToken{},
					Basic:       Basic{
						Username: "asc",
						Password: "qwerty",
					},
				},
				Host:     "10.2.3.4",
				Port:     -1,
			},
			want: false,
		},
		{
			name: "test4 over max port value",
			user: ConfluenceServerUserV1{
				Authentication: Auth{
					BearerToken: BearerToken{},
					Basic:       Basic{
						Username: "asc",
						Password: "qwerty",
					},
				},
				Host:     "10.2.3.4",
				Port:     999999999999,
			},
			want: false,
		},
		{
			name: "test5 valid entry",
			user: ConfluenceServerUserV1{
				Authentication: Auth{
					BearerToken: BearerToken{},
					Basic:       Basic{
						Username: "asc",
						Password: "qwerty",
					},
				},
				Host:     "10.2.3.4",
				Port:     8090,
			},
			want: true,
		},
		{
			name: "test6 invalid due to empty authentication struct",
			user: ConfluenceServerUserV1{
				Authentication: Auth{
					BearerToken: BearerToken{},
					Basic:       Basic{},
				},
				Host:     "10.2.3.4",
				Port:     8090,
			},
			want: false,
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
