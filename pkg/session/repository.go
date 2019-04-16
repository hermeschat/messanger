package session

import "time"

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
