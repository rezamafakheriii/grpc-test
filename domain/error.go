package domain

import (
	"grpc-test/lib"
	"grpc-test/proto"
)

func ProductNotFoundErr() error {
	return lib.ErrNotFound().WithMessage("Product not found")
}

func ErrNotEnoughCredit() error {
	return lib.ErrBadRequest().WithMessage("Not enough credit").WithGRPCErr(&proto.NotEnoughCharge{})
}
