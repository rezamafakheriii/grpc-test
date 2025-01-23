package domain

import (
	"grpc-test/lib"
)

func ProductNotFoundErr() error {
	return lib.ErrNotFound().WithMessage("Product not found")
}

func NotEnoughCreditErr() error {
	return lib.ErrBadRequest().WithMessage("Not enough credit")
}
