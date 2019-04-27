package session

import (
	"git.raad.cloud/cloud/hermes/pkg/api"
	"git.raad.cloud/cloud/hermes/pkg/repository/session"
	"github.com/pkg/errors"
)

//wtf you think it would do ? it will create session dumbass
func Create(req *api.CreateSessionRequest) *api.Response {
	//create session in mongo
	s := &session.Session{
		Node:      req.Node,
		UserAgent: req.UserAgent,
		UserID:    req.UserID,
		UserIP:    req.UserIP,
	}
	if err := session.Add(s); err != nil {
		return &api.Response{
			Code:  "500",
			Error: errors.Wrap(err, "Some error while inserting to database").Error(),
		}
	}
	return &api.Response{
		Code:  "200",
		Error: "",
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
