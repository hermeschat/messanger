package mongo

type Messagerepository struct {
	*repository
}

func NewMessagerepository() *Messagerepository {
	return &Messagerepository{
		&repository{"messages"},
	}
}

func (m *Messagerepository) Name() string {
	return m.repository.Name()
}

type Channelrepository struct {
	*repository
}

func (c *Channelrepository) Name() string {
	return c.repository.Name()
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

func (m *Userrepository) Name() string {
	return m.repository.Name()
}
