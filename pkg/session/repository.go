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

var keys = []string{"session_id", "last_activity", "user_id", "user_ip", "user_agent", "client_version", "node"}
var 
//ConstructFromMap creates a new session struct from given map
func ConstructFromMap(input map[string]interface{}) (*Session, error) {
	session := &Session{}
	err := mapstructure.Decode(input, session)
	if err != nil {
		return nil, errors.Wrap(err, "can't construct session from given map")
	}
	return session, nil
}

//ToMap generates a map from the reciever session
func (s *Session) ToMap() (map[string]interface{}, error) {
	output := map[string]interface{}{}
	err := mapstructure.Decode(s, output)
	if err != nil {
		return nil, errors.Wrap(err, "can't create a map from this session")
	}
	return output, nil
}

