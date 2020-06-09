package credentials

//Record - Aerospike configuration + credentials data
type Record struct {
	Metadata    Metadata    `json:"Metadata"`
	Owner       Owner       `json:"Owner"`
	Credentials Credentials `json:"Credentials"`
}

//Metadata - Record metadata
type Metadata struct {
	PrimaryKey string `json:"PrimaryKey"`
	LastUpdate string `json:"LastUpdate"`
	CreateTime string `json:"CreateTime"`
}

//Owner - Creation account details for grouping/fetch
type Owner struct {
	Email string `json:"Email"`
}

//Credentials - Credentials for various graph and storage services
type Credentials struct {
	GrafanaUsers map[string]DBGrafanaUser `json:"GrafanaUsers"`
}

//DBGrafanaUser - Database entry for a GrafanaUser
type DBGrafanaUser struct {
	APIKey      string `json:"ApiKey"`
	Description string `json:"Description"`
}
