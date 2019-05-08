package auth

type (
	// HTTPError custom http error to handle custom errors
	HTTPError struct {
		error      `json:"-"`
		StatusCode int    `json:"-"`
		Message    string `json:"message"`
	}
	// UnauthorizedError to handle unauthorized errors
	UnauthorizedError struct {
		HTTPError
	}
	// ForbiddenError to handle forbidden errors
	ForbiddenError struct {
		HTTPError
	}
)

// GetUnAuthorizedError returns not error associated with HTTP Unauthorized error
func GetUnAuthorizedError(message string) error {
	if message == "" {
		message = "شما به این قسمت دسترسی ندارید"
	}
	err := HTTPError{}
	err.StatusCode = 401
	err.Message = message
	return err
}

// GetForbiddenError returns not error associated with HTTP Forbidden error
func GetForbiddenError(message string) error {
	if message == "" {
		message = "شما به این قسمت دسترسی ندارید"
	}
	err := HTTPError{}
	err.StatusCode = 403
	err.Message = message
	return err
}
