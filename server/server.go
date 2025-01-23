package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"runtime"
	"runtime/debug"

	"grpc-test/domain"
	pb "grpc-test/proto"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ProductOrderService struct {
	pb.UnimplementedProductOrderServiceServer
}

func (s *ProductOrderService) ListProducts(ctx context.Context, req *pb.Empty) (*pb.ProductList, error) {
	// Simulating a product list response
	return &pb.ProductList{
		Products: []*pb.Product{
			{Id: "1", Name: "Product 1", Price: 10.0},
			{Id: "2", Name: "Product 2", Price: 20.0},
		},
	}, nil
}

func (s *ProductOrderService) PlaceOrder(ctx context.Context, req *pb.OrderRequest) (*pb.OrderResponse, error) {
	// st := status.New(codes.Unknown, "some unknown error occured")
	// return nil, st.Err()

	return &pb.OrderResponse{
		OrderId:    "12345",
		TotalPrice: 30.0,
	}, nil

	// return nil, &domain.ValidationErr{
	// 	DomainErr: &domain.DomainErr{
	// 		Code:    "P_1112222200",
	// 		Message: "some argument are invalid",
	// 	},
	// 	InvalidFields: []domain.InvalidField{
	// 		{
	// 			Field:       "order_id",
	// 			Description: "order id must be a numeric number",
	// 		},
	// 	},
	// }
}

func main() {
	serviceName := "ProductOrderService"
	debugMode := true

	server := grpc.NewServer(
		grpc.UnaryInterceptor(UnaryServerInterceptor(serviceName, debugMode)),
	)

	// Register the ProductOrderService
	pb.RegisterProductOrderServiceServer(server, &ProductOrderService{})

	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	log.Println("Starting gRPC server on portt 50051...")
	if err := server.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

func UnaryServerInterceptor(serviceName string, debugMode bool) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp interface{}, err error) {
		// Defer recovery from panic
		defer func() {
			if r := recover(); r != nil {
				if debugMode {
					log.Printf("Recovered from panic in %s: %v\n%s", info.FullMethod, r, string(debug.Stack()))
				} else {
					log.Printf("Recovered from panic in %s: %v", info.FullMethod, r)
				}

				err = recoverFrom(serviceName, r)
			}
		}()

		resp, err = handler(ctx, req)

		if err != nil {
			if domain.IsDomainError(err) {
				return nil, MapDomainErrorToGRPC(err.(*domain.ValidationErr), serviceName)
			}
			return nil, err
		}

		return resp, nil
	}
}

func MapDomainErrorToGRPC(domainErr *domain.ValidationErr, serviceName string) error {
	message := domainErr.Message

	var grpcCode codes.Code

	var notFoundErr *domain.NotFoundErr
	if errors.As(domainErr, &notFoundErr) {
		grpcCode = codes.NotFound
	}

	var validationErr *domain.ValidationErr
	if errors.As(domainErr, &validationErr) {
		grpcCode = codes.InvalidArgument
	}

	st := status.New(grpcCode, message)

	errorInfo := &errdetails.ErrorInfo{
		Domain: "product service",
		// Metadata: ,
		Reason: validationErr.Code,
	}

	stWithDetails, err := st.WithDetails(errorInfo)
	if err != nil {
		return st.Err()
	}

	if validationErr != nil {
		// If there are validation errors, include them as BadRequest details
		if len(validationErr.InvalidFields) > 0 {
			badRequest := &errdetails.BadRequest{}
			for _, ve := range validationErr.InvalidFields {
				badRequest.FieldViolations = append(badRequest.FieldViolations, &errdetails.BadRequest_FieldViolation{
					Field:       ve.Field,
					Description: ve.Description,
				})
			}
			stWithDetails, err = stWithDetails.WithDetails(badRequest)
			if err != nil {
				return st.Err()
			}
		}
	}

	return stWithDetails.Err()
}

func recoverFrom(serviceName string, r any) error {

	// Capture the stack trace
	stack := make([]byte, 64<<10) // 64 KB
	stack = stack[:runtime.Stack(stack, false)]

	st := status.New(codes.Internal, "Internal server error")

	errorInfo := &errdetails.ErrorInfo{
		Reason: "INTERNAL_SERVER_ERROR",
		Domain: serviceName,
		Metadata: map[string]string{
			"panic": fmt.Sprintf("%v", r),
			// "stack": string(stack),
		},
	}

	stWithDetails, err := st.WithDetails(errorInfo)
	if err != nil {
		return st.Err()
	}

	return stWithDetails.Err()
}
