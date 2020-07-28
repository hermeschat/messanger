package repository

type User struct {
	ID string
	Avatar string
	Handle string
}

type UserRepository interface {
	NewUser(user *User) error
	GetUser(id string) (*User, error)
}
