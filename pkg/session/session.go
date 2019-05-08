package session

import (
	"encoding/json"
	"git.raad.cloud/cloud/hermes/pkg/api"
	"git.raad.cloud/cloud/hermes/pkg/drivers/redis"
	"git.raad.cloud/cloud/hermes/pkg/repository/session"
	"github.com/pkg/errors"
	"time"
	"go.mongodb.org/mongo-driver/mongo"
)

//wtf you think it would do ? it will create session dumbass
func Create(s *session.Session) *api.Response {
	//create session in mongo
	if err := session.Add(s); err != nil {
		return &api.Response{
			Code:  "500",
			Error: errors.Wrap(err, "Some error while inserting to database").Error(),
		}
	}

	conn, err := redis.ConnectRedis()
	if err != nil {
		panic(err)
	}
	//TODO: initialize sessionID with mongo objectid
	jsonSession, err := json.Marshal(s)
	if err != nil {
		return &api.Response{
			Code:  "500",
			Error: errors.Wrap(err, "error in marshalling session to json").Error(),
		}
	}
	status := conn.Set(s.SessionID, jsonSession, time.Hour*12)
	if status.Err() != nil {
		panic(status.Err())
	}

	return &api.Response{
		Code:  "200",
		Error: "",
	}
}

func GetOrCreate(req *session.Session) *api.GetOrCreateSessionResponse {
	sess, err := GetSession(req.SessionID)
	if err != nil {
		return &api.GetOrCreateSessionResponse{
			Code: "500",
		}
	}
	if sess == nil {
		return Create(req)
	}
	return &api.GetOrCreateSessionResponse{
		Node:sess.Node,
		UserID:sess.UserID,
		ClientVersion:sess.ClientVersion,
	}
}

//it removes a session dump ass
func Destroy(req *api.DestroySessionRequest) *api.Response {
	err := session.Delete(req.SessionId)
	if err != nil {
		return &api.Response{
			Code:  "500",
			Error: errors.Wrap(err, "error while removing session").Error(),
		}
	}
	return &api.Response{
		Code:  "200",
		Error: "",
	}

}

func GetSession(sessionID string) (*session.Session, error) {
	conn, err := redis.ConnectRedis()
	if err != nil {
		return nil, errors.Wrap(err, "error in redis connection")
	}
	res, err := conn.Get(sessionID).Result()
	if err == redis.Nil {
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
