package credentials

//Credentials - Set of credentials to be validated
type Credentials struct{
	GrafanaReadUsers []GrafanaReadUser
	ConfluenceServerUsers []ConfluenceServerUser
}

//GrafanaUser - Grafana user with read access
type GrafanaReadUser struct{
	APIKey string
}

//GrafanaUser - Grafana user with read access
type ConfluenceServerUser struct{
	Username string
	Password string
}

//CredentialsCheck - Check credentials result
type CredentialsCheck struct{
	GrafanaReadUserCheck []GrafanaReadUserCheck
	ConfluenceServerUserCheck []ConfluenceUserCheck
}

//GrafanaReadUserCheck - Grafana read user check result
type GrafanaReadUserCheck struct{
	Result bool
	Cause string `json:"Cause,omitempty"`
	GrafanaReadUser
}

//ConfluenceUserCheck - Confluence write user check result
type ConfluenceUserCheck struct{
	Result bool
	Cause string `json:"Cause,omitempty"`
	ConfluenceServerUser
}