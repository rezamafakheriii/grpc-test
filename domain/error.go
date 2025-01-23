package domain

import (
	"grpc-test/lib"
)

func ProductNotFoundErr() error {
	return lib.ErrNotFound().WithMessage("Product not found")
}
