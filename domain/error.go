package domain

import (
	"grpc-test/lib"
	"grpc-test/proto"
)

func ProductNotFoundErr() error {
	return lib.ErrNotFound().WithMessage("Product not found")
}

func ErrNotEnoughCredit() error {
	return lib.ErrBadRequest().WithMessage("Not enough credit").WithProtobufError(&proto.ErrNotEnoughCharge{})
}

func ErrGatewayNotReachable() error {
	return lib.ErrBadRequest().WithMessage("Gateway not reachable").WithProtobufError(&proto.ErrGatewayNotReachable{})
}
