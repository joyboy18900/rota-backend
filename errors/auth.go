package errors

// Common error codes for authentication
const (
	ErrInvalidCredentials = "AUTH_001"
	ErrEmailExists        = "AUTH_002"
	ErrUsernameExists     = "AUTH_003"
	ErrInvalidToken       = "AUTH_004"
	ErrTokenExpired       = "AUTH_005"
	ErrTokenBlacklisted   = "AUTH_006"
)

// NewAuthError creates a new authentication error
func NewAuthError(code string, message string, err error) *Error {
	return NewError(code, message, err)
}
