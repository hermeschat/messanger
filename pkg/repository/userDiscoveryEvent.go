package repository

//UserDiscoveryEvent is the message we send to discovery channel to tell a user
//to subscribe to a certain channel in async way
type UserDiscoveryEvent struct {
	ChannelID string
	UserID    string
}
