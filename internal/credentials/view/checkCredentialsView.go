package view

type CheckCredentialsInputView struct {
	GrafanaUsers []CheckGrafanaUser `json:"GrafanaUsers"`
}

type CheckGrafanaUser struct {
	URL    string `json:"URL"`
	APIKey string `json:"APIKey"`
}

type CheckCredentialsResultView struct {
	Results []CheckGrafanaUserResultView `json:"Results"`
}

type CheckGrafanaUserResultView struct {
	URL               string `json:"URL"`
	APIKey            string `json:"APIKey"`
	MeetsRequirements bool   `json:"MeetsRequirements"`
	Error             string `json:"Error"`
}
