package credentials

import (
	"github.com/sajeevany/graph-snapper/internal/confluence"
	"github.com/sajeevany/graph-snapper/internal/grafana"
	"github.com/sirupsen/logrus"
)

func authGrafanaUsers(logger *logrus.Logger, users []CheckUserV1) ([]CheckUserResultV1, error) {

	results := make([]CheckUserResultV1, len(users))

	logger.Debug("authenticating grafana users")
	for i, user := range users {

		res, err := authenticateGrafanaUser(logger, user)
		if err != nil {
			return results, err
		}

		results[i] = res
	}
	logger.Debug("done authenticating grafana users")


	return results, nil
}

func authenticateGrafanaUser(logger *logrus.Logger, gu CheckUserV1) (CheckUserResultV1, error) {

	isValid, rErr := grafana.IsValidLogin(logger, gu.Auth, gu.Host, gu.Port)
	if rErr != nil {
		logger.WithFields(gu.GetFields()).Errorf("Error checking if grafana user has login access. <%v>", rErr)
		return CheckUserResultV1{}, rErr
	}

	if isValid {
		return CheckUserResultV1{
			Result:      true,
			CheckUserV1: gu,
		}, nil
	} else {
		return CheckUserResultV1{
			Result:      false,
			Cause:       "Unauthorized. Received 401.",
			CheckUserV1: gu,
		}, nil
	}
}

func authConfluenceUsers(logger *logrus.Logger, users []CheckUserV1) ([]CheckUserResultV1, error) {

	results := make([]CheckUserResultV1, len(users))

	logger.Debug("authenticating confluence server users")
	for i, user := range users {
		res, err := authenticateConfluenceUser(logger, user)
		if err != nil {
			return results, err
		}

		results[i] = res
	}
	logger.Debug("done authenticating confluence server users")


	return results, nil
}

func authenticateConfluenceUser(logger *logrus.Logger, cu CheckUserV1) (CheckUserResultV1, error) {

	hasWriteAccess, rErr := confluence.HasWriteAccess(logger, cu.Host, cu.Port, cu.Auth)
	if rErr != nil {
		logger.WithFields(cu.GetFields()).Errorf("Error checking if confluence user has write access. <%v>", rErr)
		return CheckUserResultV1{}, rErr
	}

	//Check if user has write access
	if hasWriteAccess {
		logger.Debug("User has write access")
		return CheckUserResultV1{
			Result:      true,
			CheckUserV1: cu,
		}, nil
	} else {
		//Confluence returns
		logger.Debug("User does not have write access")
		return CheckUserResultV1{
			Result:      false,
			Cause:       "User does not have access mode READ_WRITE",
			CheckUserV1: cu,
		}, nil
	}
}
