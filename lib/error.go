package lib

import (
	"github.com/revotech-group/go-lib/errors"
)

const (
	NameServiceUnavailable  = "ServiceUnavailable"
	NameTooManyRequests     = "TooManyRequestsError"
	NameBadRequest          = "BadRequestError"
	NameInternalServerError = "InternalServerError"
	NameNotFound            = "NotFoundError"
	NameForbidden           = "ForbiddenError"
	NameAlreadyExists       = "AlreadyExistsError"
	NameUnauthorizedAccess  = "UnauthorizedAccessError"
)

// #TODO: we can pass code as parameter in these functions so it can be customized in caller side with meaningful code
func ErrNotFound() errors.AppError {
	return errors.NewAppError(NameNotFound, "Not found or does not exist", 404)
}

func ErrInternalServerError() errors.AppError {
	return errors.NewAppError(NameInternalServerError, "Internal server error", 500)
}

func ErrBadRequest() errors.AppError {
	return errors.NewAppError(NameBadRequest, "Bad request, invalid or missing parameter", 400)
}
