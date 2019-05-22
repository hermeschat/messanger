package session

import (
	"encoding/json"
	"git.raad.cloud/cloud/hermes/pkg/drivers/redis"
	"git.raad.cloud/cloud/hermes/pkg/repository/session"
	"github.com/pkg/errors"
	"github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)


type CreateSession struct {
	UserID string
	UserIP string
	ClientVersion string
	Node string
}

//var State = &map[string]*stan.Conn{}

//wtf you think it would do ? it will create session dumbass
func Create(cs *CreateSession) (*session.Session,error) {
	sid, err :=  uuid.NewV4()
	if err != nil {
		return nil, errors.Wrap(err, "cant create new uuid")
	}
	s := &session.Session{
		UserID:cs.UserID,
		UserIP:cs.UserIP,
		ClientVersion:cs.ClientVersion,
		Node:cs.Node,
		SessionID: sid.String(),
	}
	//create session in mongo
	if err := session.Add(s); err != nil {
		return nil, errors.Wrap(err, "error in creating")
	}
	logrus.Println("added to mongo db")
	conn, err := redis.ConnectRedis()
	if err != nil {
		return nil, errors.Wrap(err, "error in connecting to redis")
	}
	//TODO: initialize sessionID with mongo objectid
	jsonSession, err := json.Marshal(s)
	if err != nil {
		return nil, errors.Wrap(err, "error in marshaling json")
	}
	status := conn.Set(s.SessionID, jsonSession, time.Hour*12)
	if status.Err() != nil {
		logrus.Errorf("could not set redis key :%s", err.Error())
		return s, nil
	}
	return s, nil
}

func GetOrCreate(sessionID string,cs *CreateSession) (*session.Session, error) {
	sess, err := GetSession(sessionID)
	if err != nil {
		return nil, errors.Wrap(err, "error in getting session")
	}
	if sess == nil {
		return Create(cs)
	}
	return sess, nil
}

//it removes a session dump ass
func Destroy(sessionID string) error {
	err := session.Delete(sessionID)
	if err != nil {
		return errors.Wrap(err, "error in destroying session")
	}
	return errors.Wrap(err, "error in deleting session")

}

func GetSession(sessionID string) (*session.Session, error) {
	conn, err := redis.ConnectRedis()
	if err != nil {
		return nil, errors.Wrap(err, "error in redis connection")
	}
	res, err := conn.Get(sessionID).Result()
	if err == redis.Nil {
		logrus.Infof("session not found in redis trying db")
		return GetSessionFromDB(sessionID)
	}
	if err != nil {
		return nil, errors.Wrap(err, "error in getting redis key")
	}
	s := &session.Session{}
	err = json.Unmarshal([]byte(res), s)
	if err != nil {
		return nil, errors.Wrap(err, "cant unmarshall session")
	}
	return s, nil
}

func GetSessionFromDB(sessionID string) (*session.Session, error) {
	s, err := session.Get(sessionID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, errors.Wrap(err, "Error in get session from mongo")
	}
	return s, nil
}
