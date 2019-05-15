package session

import (
	"time"

	"git.raad.cloud/cloud/hermes/pkg/drivers/mongo"

	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
)

//Session ...
type Session struct {
	SessionID     string
	LastActivity  time.Time
	UserID        string
	UserIP        string
	UserAgent     string
	ClientVersion string
	Node          string
}

//ConstructFromMap ...
func ConstructFromMap(m map[string]interface{}) (*Session, error) {
	user := &Session{}
	err := mapstructure.Decode(m, user)
	if err != nil {
		return nil, errors.Wrap(err, "can't construct user from given map")
	}
	return user, nil
}

func Get(id string) (*Session, error) {
	s, err := mongo.FindOneById("sessions", id)
	if err != nil {
		return nil, errors.Wrap(err, "can't find session with given id")
	}
	session := &Session{}
	err = mapstructure.Decode(s, session)
	if err != nil {
		return nil, errors.Wrap(err, "can't construct session from given map from mongo")
	}
	return session, nil
}

func (s *Session) ToMap() (map[string]interface{}, error) {
	m := map[string]interface{}{}
	err := mapstructure.Decode(s, &m)
	if err != nil {
		return nil, errors.Wrap(err, "can't create map from this session")
	}
	return m, nil
}
func GetAll(query map[string]interface{}) (*[]*Session, error) {
	s, err := mongo.FindAll("sessions", query)
	if err != nil {
		return nil, errors.Wrap(err, "can't find session with given query")
	}
	sessions := &[]*Session{}
	err = mapstructure.Decode(s, sessions)
	if err != nil {
		return nil, errors.Wrap(err, "can't construct session from given map from mongo")
	}
	return sessions, nil
}

//Add adds session
func Add(session *Session) error {

	err := mongo.InsertOne("sessions", session)
	if err != nil {
		return errors.Wrap(err, "can't add this session to mongodb")
	}
	return nil
}

func AddAll(sessions []Session) error {
	sessionsMap := []interface{}{}
	for _, sess := range sessions {
		m := map[string]interface{}{}
		err := mapstructure.Decode(sess, m)
		if err != nil {
			return errors.Wrap(err, "can't convert sessions to map")
		}
		sessionsMap = append(sessionsMap, m)
	}
	err := mongo.InsertAll("sessions", sessionsMap)

	if err != nil {
		return errors.Wrap(err, "can't add this session to mongodb")
	}
	return nil

}
func Delete(id string) error {
	err := mongo.DeleteById("sessions", id)
	if err != nil {
		return errors.Wrap(err, "can't delete this session from mongodb")
	}
	return nil
}

func DeleteAll(query map[string]interface{}) error {
	err := mongo.DeleteAllMatched("sessions", query)
	if err != nil {
		return errors.Wrap(err, "can't delete with given query from mongo")
	}
	return nil
}

//TODO add count of updated docs
func UpdateOne(id string, query map[string]interface{}) error {
	err := mongo.UpdateOne("sessions", id, query)
	if err != nil {
		return errors.Wrap(err, "can't update this id with given query")
	}
	return nil
}

func UpdateAll(selector map[string]interface{}, updator map[string]interface{}) error {
	err := mongo.UpdateAllMatched("sessions", selector, updator)
	if err != nil {
		return errors.Wrap(err, "can't update session with given query")
	}
	return nil
}
