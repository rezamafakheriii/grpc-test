package domain

import "fmt"

// DomainErr is the base error type.
type DomainErr struct {
	Code    string
	Message string
}

func (ge *DomainErr) Error() string {
	return fmt.Sprintf("err code: %s, err message: %s", ge.Code, ge.Message)
}

// ValidationErr is a derived error type that embeds DomainErr.
type ValidationErr struct {
	*DomainErr
	InvalidFields []InvalidField
}

type InvalidField struct {
	Field       string
	Description string
}

// NotFoundErr is another derived error type that embeds DomainErr.
type NotFoundErr struct {
	*DomainErr
	Resource string
}

func ProductNotFoundError(productID string) *NotFoundErr {
	return &NotFoundErr{
		Resource: "product",
		DomainErr: &DomainErr{
			Message: fmt.Sprintf("product %s not found.", productID),
			Code:    "P_1122300",
		},
	}
}

// IsDomainError checks if the error is of type DomainErr or any type that embeds DomainErr.
func IsDomainError(err error) bool {
	if err == nil {
		return false
	}

	// Check if the error is of type *DomainErr
	if _, ok := err.(*DomainErr); ok {
		return true
	}

	// Check if the error is of type *ValidationErr
	if _, ok := err.(*ValidationErr); ok {
		return true
	}

	// Check if the error is of type *NotFoundErr
	if _, ok := err.(*NotFoundErr); ok {
		return true
	}

	return false
}
