package credentials

import (
	"github.com/sajeevany/graph-snapper/internal/confluence"
	"github.com/sajeevany/graph-snapper/internal/grafana"
	"github.com/sirupsen/logrus"
)

func authGrafanaUsers(logger *logrus.Logger, users []CheckGrafanaReadUserV1) ([]CheckGrafanaReadUserResultV1, error) {

	results := make([]CheckGrafanaReadUserResultV1, len(users))

	logger.Debug("authenticating users")
	for i, user := range users {

		res, err := authenticateGrafanaUser(logger, user)
		if err != nil {
			return results, err
		}

		results[i] = res
	}

	return results, nil
}

func authenticateGrafanaUser(logger *logrus.Logger, gu CheckGrafanaReadUserV1) (CheckGrafanaReadUserResultV1, error) {

	isValid, rErr := grafana.IsValidLogin(logger, gu.APIKey, gu.Host, gu.Port)
	if rErr != nil {
		logger.WithFields(gu.GetFields()).Errorf("Error checking if grafana user has login access. <%v>", rErr)
		return CheckGrafanaReadUserResultV1{}, rErr
	}

	if isValid {
		return CheckGrafanaReadUserResultV1{
			Result:                 true,
			CheckGrafanaReadUserV1: gu,
		}, nil
	} else {
		return CheckGrafanaReadUserResultV1{
			Result:                 false,
			Cause:                  "Unauthorized. Received 401.",
			CheckGrafanaReadUserV1: gu,
		}, nil
	}
}

func authConfluenceUsers(logger *logrus.Logger, users []CheckConfluenceServerUserV1) ([]CheckConfluenceUserResultV1, error) {

	results := make([]CheckConfluenceUserResultV1, len(users))

	logger.Debug("authenticating confluence server users")
	for i, user := range users {
		res, err := authenticateConfluenceUser(logger, user)
		if err != nil {
			return results, err
		}

		results[i] = res
	}

	return results, nil
}

func authenticateConfluenceUser(logger *logrus.Logger, cu CheckConfluenceServerUserV1) (CheckConfluenceUserResultV1, error) {

	hasWriteAccess, rErr := confluence.HasWriteAccess(logger, cu.Host, cu.Port, cu.Username, cu.Password)
	if rErr != nil {
		logger.WithFields(cu.GetFields()).Errorf("Error checking if confluence user has write access. <%v>", rErr)
		return CheckConfluenceUserResultV1{}, rErr
	}

	//Check if user has write access
	if hasWriteAccess {
		logger.Debug("User has write access")
		return CheckConfluenceUserResultV1{
			Result:                      true,
			CheckConfluenceServerUserV1: cu,
		}, nil
	} else {
		//Confluence returns
		logger.Debug("User does not have write access")
		return CheckConfluenceUserResultV1{
			Result:                      false,
			Cause:                       "User does not have access mode READ_WRITE",
			CheckConfluenceServerUserV1: cu,
		}, nil
	}
}
