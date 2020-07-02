package handler

type CheckCredentialsResult struct {
	Results []CheckGrafanaUserResult
}

type CheckGrafanaUserResult struct {
	URL               string
	APIKey            string
	MeetsRequirements bool
	Error             string
}

