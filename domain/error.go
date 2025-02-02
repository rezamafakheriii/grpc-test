package domain

import (
	"grpc-test/lib"
	"grpc-test/proto"

	"github.com/revotech-group/go-lib/errors"
)

func ProductNotFoundErr() error {
	return lib.ErrNotFound().WithMessage("Product not found")
}

func ErrNotEnoughCredit() error {
	return lib.ErrBadRequest().WithMessage("Not enough credit").WithProtobufError(&proto.ErrNotEnoughCharge{})
}

func ErrGatewayNotReachable() error {
	return errors.NewAppError(lib.NameTooManyRequests, "Bad request, invalid or missing parameter", 500).WithMessage("Gateway not reachable").WithProtobufError(&proto.ErrGatewayNotReachable{})
}
