package main

import (
	"context"
	"log"
	"log/slog"
	"net"
	"runtime/debug"

	"grpc-test/domain"
	pb "grpc-test/proto"

	errlib "github.com/revotech-group/go-lib/errors"
	logger "github.com/revotech-group/go-lib/log"
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
	// return nil, fmt.Errorf("some unknown error")

	// return &pb.OrderResponse{
	// 	OrderId:    "12345",
	// 	TotalPrice: 30.0,
	// }, nil

	// panic("some error occured")

	return nil, domain.ProductNotFoundErr()
}

func main() {
	serviceName := "ProductOrderService"
	debugMode := true
	logger.SetupDefaultLogger(slog.LevelDebug, true)

	server := grpc.NewServer(
		grpc.UnaryInterceptor(UnaryServerInterceptor(serviceName, debugMode)),
	)

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

			if IsAppError(err) {
				appErr := err.(errlib.AppError)

				slog.Error("app error",
					slog.String("method", info.FullMethod),
					slog.String("error", appErr.Error()),
					slog.Any("stack_trace", appErr.StackTrace()),
				)

				return nil, MapAppErrorToGRPC(err.(errlib.AppError), serviceName)
			}

			slog.Error("gRPC error",
				slog.String("method", info.FullMethod),
				slog.String("error", err.Error()),
			)
			return nil, err
		}

		return resp, nil
	}
}

func MapAppErrorToGRPC(appErr errlib.AppError, serviceName string) error {
	message := appErr.GetMessage()

	var grpcCode codes.Code

	switch appErr.GetName() {
	case errlib.NameNotFound:
		grpcCode = codes.NotFound
	case errlib.NameBadRequest:
		grpcCode = codes.InvalidArgument
	case errlib.NameInternalServerError:
		grpcCode = codes.Internal
	default:
		grpcCode = codes.Unknown
	}

	st := status.New(grpcCode, message)

	errorInfo := &errdetails.ErrorInfo{
		Domain: serviceName,
	}

	stWithDetails, err := st.WithDetails(errorInfo)
	if err != nil {
		return st.Err()
	}

	return stWithDetails.Err()
}

func recoverFrom(serviceName string, _ any) error {

	// // Capture the stack trace
	// stack := make([]byte, 64<<10) // 64 KB
	// stack = stack[:runtime.Stack(stack, false)]

	st := status.New(codes.Internal, "Internal server error")

	errorInfo := &errdetails.ErrorInfo{
		Reason: "INTERNAL_SERVER_ERROR",
		Domain: serviceName,
	}

	stWithDetails, err := st.WithDetails(errorInfo)
	if err != nil {
		return st.Err()
	}

	return stWithDetails.Err()
}

func IsAppError(err error) bool {
	if _, ok := err.(errlib.AppError); ok {
		return true
	}
	return false
}
