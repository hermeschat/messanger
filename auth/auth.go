package auth

type ID string
type Token string
//AuthService....
type Service interface {
	NewUser() ID
	GetUser(id ID) error
	GetToken(id ID) (Token, error)
	VerifyToken(token Token) error
}
