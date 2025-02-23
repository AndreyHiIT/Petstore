package errors

import "fmt"

var (
	ErrTokenType        = fmt.Errorf("token type unknown")
	ErrTokenExtractUser = fmt.Errorf("type assertion to user err")
	ErrPetNotFound      = fmt.Errorf("no pet found with the provided id")
)
