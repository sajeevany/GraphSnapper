package credentials

import (
	"github.com/sajeevany/graph-snapper/internal/grafana"
	"github.com/sirupsen/logrus"
)

func authGrafanaUsers(logger *logrus.Logger, users []GrafanaReadUser) ([]GrafanaReadUserCheck, error) {

	results := make([]GrafanaReadUserCheck, len(users))

	logger.Debug("authenticating users")
	for i, user := range users {
		results[i] = authenticateGrafanaUser(logger, user)
	}

	return results, nil
}

func authenticateGrafanaUser(logger *logrus.Logger, gu GrafanaReadUser) GrafanaReadUserCheck {

	isValid, rErr := grafana.IsValidLogin(logger, gu.APIKey, gu.Host, gu.Port)
	if rErr != nil {
		return GrafanaReadUserCheck{
			Result:          false,
			Cause:           rErr.Error(),
			GrafanaReadUser: gu,
		}
	}

	if isValid {
		return GrafanaReadUserCheck{
			Result:          true,
			GrafanaReadUser: gu,
		}
	}

	return GrafanaReadUserCheck{
		Result:          false,
		Cause:           "Not validated",
		GrafanaReadUser: gu,
	}
}
