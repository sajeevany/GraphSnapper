package credentials

import (
	"github.com/sajeevany/graph-snapper/internal/confluence"
	"github.com/sajeevany/graph-snapper/internal/grafana"
	"github.com/sirupsen/logrus"
)

func authGrafanaUsers(logger *logrus.Logger, users []GrafanaReadUser) ([]GrafanaReadUserCheck, error) {

	results := make([]GrafanaReadUserCheck, len(users))

	logger.Debug("authenticating users")
	for i, user := range users {

		res, err := authenticateGrafanaUser(logger, user)
		if err != nil{
			return results, err
		}

		results[i] = res
	}

	return results, nil
}

func authenticateGrafanaUser(logger *logrus.Logger, gu GrafanaReadUser) (GrafanaReadUserCheck, error) {

	isValid, rErr := grafana.IsValidLogin(logger, gu.APIKey, gu.Host, gu.Port)
	if rErr != nil {
		logger.WithFields(gu.GetFields()).Errorf("Error checking if grafana user has login access. <%v>", rErr)
		return GrafanaReadUserCheck{}, rErr
	}

	if isValid {
		return GrafanaReadUserCheck{
			Result:          true,
			GrafanaReadUser: gu,
		}, nil
	}else{
		return GrafanaReadUserCheck{
			Result:          false,
			Cause:           "Not validated",
			GrafanaReadUser: gu,
		}, nil
	}
}

func authConfluenceUsers(logger *logrus.Logger, users []ConfluenceServerUser) ([]ConfluenceUserCheck, error) {

	results := make([]ConfluenceUserCheck, len(users))

	logger.Debug("authenticating confluence server users")
	for i, user := range users {
		res, err := authenticateConfluenceUser(logger, user)
		if err != nil{
			return results, err
		}

		results[i] = res
	}

	return results, nil
}

func authenticateConfluenceUser(logger *logrus.Logger, cu ConfluenceServerUser) (ConfluenceUserCheck, error) {

	hasWriteAccess, rErr := confluence.HasWriteAccess(logger, cu.Host, cu.Port, cu.Username, cu.Password)
	if rErr != nil {
		logger.WithFields(cu.GetFields()).Errorf("Error checking if confluence user has write access. <%v>", rErr)
		return ConfluenceUserCheck{}, rErr
	}

	//Check if user has write access
	if hasWriteAccess {
		logger.Debug("User has write access")
		return ConfluenceUserCheck{
			Result:          true,
			ConfluenceServerUser: cu,
		}, nil
	}else{
		//Confluence returns
		logger.Debug("User does not have write access")
		return ConfluenceUserCheck{
			Result:          false,
			Cause:           "User does not have access mode READ_WRITE",
			ConfluenceServerUser: cu,
		}, nil
	}
}