package interceptor

import (
	"context"
	"log"
	"log/slog"

	errlib "github.com/revotech-group/go-lib/errors"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/protoadapt"
)

func UnaryServerInterceptor(serviceName string, debugMode bool) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp interface{}, err error) {

		// defer func() {
		// 	if r := recover(); r != nil {
		// 		if debugMode {
		// 			slog.Error("Recovered from panic",
		// 				slog.String("method", info.FullMethod),
		// 				slog.Any("error", r),
		// 				slog.String("stack_trace", string(debug.Stack())),
		// 			)
		// 		} else {
		// 			slog.Error("Recovered from panic",
		// 				slog.String("method", info.FullMethod),
		// 				slog.Any("error", r),
		// 			)
		// 		}

		// 		err = recoverFrom(serviceName, r)
		// 	}
		// }()

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

	errorInfo := appErr.ToGRPCErrorInfo(serviceName)

	stWithDetails, err := st.WithDetails(errorInfo)
	if err != nil {
		return st.Err()
	}

	// if appErr.GetGRPCErr() != nil {
	// 	log.Printf("grpc err detected")
	// 	stWithDetails, err := stWithDetails.WithDetails(appErr.GetGRPCErr())
	// 	if err != nil {
	// 		return st.Err()
	// 	}
	// 	return stWithDetails.Err()
	// }

	if grpcMsg, ok := appErr.(protoadapt.MessageV1); ok {
		stWithDetails, err := stWithDetails.WithDetails(grpcMsg)
		if err != nil {
			return st.Err()
		}
		return stWithDetails.Err()
	} else {
		log.Printf("not campatible")
	}

	return stWithDetails.Err()
}

func recoverFrom(serviceName string, _ any) error {
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
