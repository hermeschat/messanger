package auth

// IdentityService service info of identity
type IdentityService struct {
	Permissions int64
}

// Identity for user authentication
type Identity struct {
	// NameID   string
	ID            string
	Token         string
	Roles         map[string]bool
	MerchantRoles map[string][]string
	Service       IdentityService
	AppId         string
}
