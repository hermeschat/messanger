package user

import "time"

//User is User
type User struct {
	AccountID      string `bson:"_id" json:"_id"`
	UserName       string
	Name           string
	ProfilePicture string
	MobilePhone    string
	Active         bool
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func ConstructFromMap(m map[string]interface{}) (*User, error) {

}
