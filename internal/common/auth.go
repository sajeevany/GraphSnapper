package common

import (
	"fmt"
	"github.com/sajeevany/graph-snapper/internal/logging"
	"github.com/sirupsen/logrus"
	"net/http"
)

const (
	BearerTokenAuthType = "BEARER_TOKEN"
	BasicAuthType       = "BASIC"
)

//Auth
type Auth struct {
	BearerToken BearerToken
	Basic       Basic
}

//SetAuthHeader - sets authentication header with the highest priority
func SetAuthHeader(logger *logrus.Logger, auth Auth, req *http.Request) {

	switch getHighestPriorityAuthType(auth) {
	case BearerTokenAuthType:
		logger.Debug("Request has bearer token. Setting bearer auth header.")
		setBearerAuthHeader(req, auth.BearerToken)
	case BasicAuthType:
		logger.Debug("Request has basic auth type. Setting basic auth header.")
		req.SetBasicAuth(auth.Basic.Username, auth.Basic.Password)
	default:
		logger.Info("No authentication type was found. Skipping set auth header.")
	}
}

func getHighestPriorityAuthType(auth Auth) string {

	if auth.BearerToken != (BearerToken{}) {
		return BearerTokenAuthType
	}

	if auth.Basic != (Basic{}) {
		return BasicAuthType
	}

	return ""
}

func setBearerAuthHeader(req *http.Request, token BearerToken) {
	bearer := fmt.Sprintf("Bearer %s", token.Token)
	req.Header.Add("Authorization", bearer)
}

func (a Auth) GetFields() logrus.Fields {
	return logrus.Fields{
		"Basic":       a.Basic.GetFields(),
		"BearerToken": a.BearerToken.GetFields(),
	}
}

type Basic struct {
	Username string
	Password string
}

func (a Basic) GetFields() logrus.Fields {
	return logrus.Fields{
		"Username": logging.RedactNonEmpty(a.Username),
		"Password": logging.RedactNonEmpty(a.Password),
	}
}

type BearerToken struct {
	Token string
}

func (a BearerToken) GetFields() logrus.Fields {
	return logrus.Fields{
		"Token": logging.RedactNonEmpty(a.Token),
	}
}
