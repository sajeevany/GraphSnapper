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

	//Should be identical to json name of Auth fields so that direct unmarshalling will work
	BearerTokenASName = "BearerToken"
	BasicASName = "Basic"
)

//Auth
type Auth struct {
	BearerToken BearerToken
	Basic       Basic
}

func (a Auth) GetRedactedLog() logrus.Fields {

	if a.BearerToken != (BearerToken{}) {
		return logrus.Fields{
			"BearerToken": logging.RedactNonEmpty(a.BearerToken.Token),
		}
	}

	if a.Basic != (Basic{}) {
		return logrus.Fields{
			"User":     logging.RedactNonEmpty(a.Basic.Username),
			"Password": logging.RedactNonEmpty(a.Basic.Password),
		}
	}

	return logrus.Fields{
		"AuthContents": "Empty",
	}
}

//IsValid - returns the validity check result of the highest priority auth type provided
func (a Auth) IsValid() bool{
	if a.BearerToken != (BearerToken{}) {
		return a.BearerToken.IsValid()
	}

	if a.Basic != (Basic{}) {
		return a.Basic.IsValid()
	}

	return false
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

func (a Auth) ToAerospikeBinMap() map[string]interface{} {

	authBM := make(map[string]interface{}, 2)
	authBM[BearerTokenASName] = a.BearerToken.ToAerospikeBinMap()
	authBM[BasicASName] = a.Basic.ToAerospikeBinMap()

	return authBM
}

type Basic struct {
	Username string
	Password string
}

func (b Basic) GetFields() logrus.Fields {
	return logrus.Fields{
		"Username": logging.RedactNonEmpty(b.Username),
		"Password": logging.RedactNonEmpty(b.Password),
	}
}

func (b Basic) ToAerospikeBinMap() map[string]string {
	return map[string]string{
		"Username": b.Username,
		"Password": b.Password,
	}
}

func (b Basic) IsValid() bool {
	return b.Username != "" && b.Password != ""
}

type BearerToken struct {
	Token string
}

func (a BearerToken) GetFields() logrus.Fields {
	return logrus.Fields{
		"Token": logging.RedactNonEmpty(a.Token),
	}
}

func (bt BearerToken) ToAerospikeBinMap() map[string]string {
	return map[string]string{
		"Token": bt.Token,
	}
}

func (a BearerToken) IsValid() bool {
	return a.Token != ""
}
