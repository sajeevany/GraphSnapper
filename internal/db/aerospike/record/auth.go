package record

import "github.com/sirupsen/logrus"

//DbAuth
type DBAuth struct {
	Basic       Basic
	BearerToken BearerToken
}

func (a DBAuth) GetFields() logrus.Fields {
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
		"Username": a.Username,
		"Password": a.Password,
	}
}

type BearerToken struct {
	Token string
}

func (a BearerToken) GetFields() logrus.Fields {
	return logrus.Fields{
		"Token": a.Token,
	}
}
