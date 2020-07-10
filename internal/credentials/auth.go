package credentials

import (
	"github.com/sajeevany/graph-snapper/internal/confluence"
	"github.com/sajeevany/graph-snapper/internal/grafana"
	"github.com/sirupsen/logrus"
)

func authGrafanaUsers(logger *logrus.Logger, users []CheckUserV1) []CheckUserResultV1 {

	results := make([]CheckUserResultV1, len(users))

	logger.Debug("authenticating grafana users")
	for i, user := range users {
		results[i] = authenticateGrafanaUser(logger, user)
	}
	logger.Debug("done authenticating grafana users")

	return results
}

func authenticateGrafanaUser(logger *logrus.Logger, gu CheckUserV1) CheckUserResultV1 {

	isValid, rErr := grafana.IsValidLogin(logger, gu.Auth, gu.Host, gu.Port)
	logger.Infof("Received %v %v for %+v", isValid, rErr, gu)
	if rErr != nil {
		logger.WithFields(gu.GetFields()).Errorf("Error checking if grafana user has login access. <%v>", rErr)
		return CheckUserResultV1{
			Result:      false,
			Cause:       rErr.Error(),
			CheckUserV1: gu,
		}
	}

	if isValid {
		return CheckUserResultV1{
			Result:      true,
			CheckUserV1: gu,
		}
	} else {
		return CheckUserResultV1{
			Result:      false,
			Cause:       "Unauthorized. Received 401",
			CheckUserV1: gu,
		}
	}
}

func authConfluenceUsers(logger *logrus.Logger, users []CheckUserV1) []CheckUserResultV1 {

	results := make([]CheckUserResultV1, len(users))

	logger.Debug("authenticating confluence server users")
	for i, user := range users {
		results[i] = authenticateConfluenceUser(logger, user)
	}
	logger.Debug("done authenticating confluence server users")

	return results
}

func authenticateConfluenceUser(logger *logrus.Logger, cu CheckUserV1) CheckUserResultV1 {

	hasWriteAccess, rErr := confluence.HasWriteAccess(logger, cu.Host, cu.Port, cu.Auth)
	if rErr != nil {
		logger.WithFields(cu.GetFields()).Errorf("Error checking if confluence user has write access. <%v>", rErr)
		return CheckUserResultV1{
			Result:      false,
			Cause:       rErr.Error(),
			CheckUserV1: cu,
		}
	}

	//Check if user has write access
	if hasWriteAccess {
		logger.Debug("User has write access")
		return CheckUserResultV1{
			Result:      true,
			CheckUserV1: cu,
		}
	} else {
		//Confluence returns
		logger.Debug("User does not have write access")
		return CheckUserResultV1{
			Result:      false,
			Cause:       "User does not have access mode READ_WRITE",
			CheckUserV1: cu,
		}
	}
}
