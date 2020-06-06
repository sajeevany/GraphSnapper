package user

import "fmt"

type AddUsersModel struct {
	GrafanaUsers []GrafanaUser
}

//IsValid - Returns true if valid. If in valid returns false and an error
func (am AddUsersModel) IsValid() error{

	//Return false if entry is empty
	if len(am.GrafanaUsers) == 0{
		return fmt.Errorf("GrafanaUsers array is empty")
	}

	//Validate grafana user
	for _, user := range am.GrafanaUsers{
		if !user.isValid(){
			return fmt.Errorf("invalid grafana user provided key <%v> desc <%v>", user.APIKey, user.Description)
		}
	}

	return nil
}

type GrafanaUser struct {
	APIKey string
	Description string
}

func (gu GrafanaUser) isValid() bool{
	return gu.APIKey != ""
}

type StoredUsers struct {
	Users []Users
}

type Users struct {
	Key string
	Description string
}