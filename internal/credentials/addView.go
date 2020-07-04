package credentials

//AddCredentialsV1 - Request to add users to account with specified ID
type AddCredentialsV1 struct {
	AccountID             string
	GrafanaReadUsers      []AddGrafanaReadUserV1      `json:"GrafanaReadUsers"`
	ConfluenceServerUsers []AddConfluenceServerUserV1 `json:"ConfluenceServerUsers"`
}

//AddedCredentials - Users added to specified account
type AddedCredentialsV1 struct {
	GrafanaReadUsers      map[string]AddGrafanaReadUserV1      `json:"GrafanaReadUsers"`
	ConfluenceServerUsers map[string]AddConfluenceServerUserV1 `json:"ConfluenceServerUsers"`
}

//AddGrafanaReadUserV1 - Grafana user with read access
type AddGrafanaReadUserV1 struct {
	APIKey      string
	Host        string
	Port        int
	Description string
}

//AddConfluenceServerUserV1 - confluence user with write access
type AddConfluenceServerUserV1 struct {
	Username string
	Password string
	Host     string
	Port     int
}
