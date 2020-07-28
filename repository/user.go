package repository

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID primitive.ObjectID
	Avatar string
	Handle string
}

type UserRepository interface {
	NewUser(user *User) error
	GetUser(id primitive.ObjectID) (*User, error)
}
