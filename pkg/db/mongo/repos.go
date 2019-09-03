package mongo

type Messagerepository struct {
	*repository
}

func NewMessagerepository() *Messagerepository {
	return &Messagerepository{
		&repository{"messages"},
	}
}

func (*Messagerepository) Name() string {
	return "messages"
}

type Channelrepository struct {
	*repository
}

func NewChannelrepository() *Channelrepository {
	return &Channelrepository{repository: &repository{"channels"}}
}

type Userrepository struct {
	*repository
}

func NewUserrepository() *Userrepository {
	return &Userrepository{repository: &repository{"users"}}
}
