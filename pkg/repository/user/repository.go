package user

import (
	"time"

	"git.raad.cloud/cloud/hermes/pkg/drivers/mongo"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
)

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

//ConstructFromMap ...
func ConstructFromMap(m map[string]interface{}) (*User, error) {
	user := &User{}
	err := mapstructure.Decode(m, user)
	if err != nil {
		return nil, errors.Wrap(err, "can't construct user from given map")
	}
	return user, nil
}

func Get(id string) (*User, error) {
	s, err := mongo.FindOneById("users", id)
	if err != nil {
		return nil, errors.Wrap(err, "can't find user with given id")
	}
	user := &User{}
	err = mapstructure.Decode(s, user)
	if err != nil {
		return nil, errors.Wrap(err, "can't construct user from given map from mongo")
	}
	return user, nil
}

func (s *User) ToMap() (map[string]interface{}, error) {
	m := map[string]interface{}{}
	err := mapstructure.Decode(s, m)
	if err != nil {
		return nil, errors.Wrap(err, "can't create map from this user")
	}
	return m, nil
}
func GetAll(query map[string]interface{}) (*[]*User, error) {
	s, err := mongo.FindAll("users", query)
	if err != nil {
		return nil, errors.Wrap(err, "can't find user with given query")
	}
	users := &[]*User{}
	err = mapstructure.Decode(s, users)
	if err != nil {
		return nil, errors.Wrap(err, "can't construct user from given map from mongo")
	}
	return users, nil
}

//Add adds user
func Add(user *User) error {

	err := mongo.InsertOne("users", user)
	if err != nil {
		return errors.Wrap(err, "can't add this user to mongodb")
	}
	return nil
}

func AddAll(users []User) error {
	usersMap := []interface{}{}
	for _, sess := range users {
		m := map[string]interface{}{}
		err := mapstructure.Decode(sess, m)
		if err != nil {
			return errors.Wrap(err, "can't convert users to map")
		}
		usersMap = append(usersMap, m)
	}
	err := mongo.InsertAll("users", usersMap)

	if err != nil {
		return errors.Wrap(err, "can't add this user to mongodb")
	}
	return nil

}
func Delete(id string) error {
	err := mongo.DeleteById("users", id)
	if err != nil {
		return errors.Wrap(err, "can't delete this user from mongodb")
	}
	return nil
}

func DeleteAll(query map[string]interface{}) error {
	err := mongo.DeleteAllMatched("users", query)
	if err != nil {
		return errors.Wrap(err, "can't delete with given query from mongo")
	}
	return nil
}

//TODO add count of updated docs
func UpdateOne(id string, query map[string]interface{}) error {
	err := mongo.UpdateOne("users", id, query)
	if err != nil {
		return errors.Wrap(err, "can't update this id with given query")
	}
	return nil
}

func UpdateAll(selector map[string]interface{}, updator map[string]interface{}) error {
	err := mongo.UpdateAllMatched("users", selector, updator)
	if err != nil {
		return errors.Wrap(err, "can't update user with given query")
	}
	return nil
}
