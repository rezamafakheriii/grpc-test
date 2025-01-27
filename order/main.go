package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net"
	"reflect"

	"grpc-test/interceptor"
	"grpc-test/proto"
	pb "grpc-test/proto" // Replace with the correct import path

	"github.com/revotech-group/go-lib/errors"

	logger "github.com/revotech-group/go-lib/log"
	"google.golang.org/grpc"
)

type orderServer struct {
	pb.UnimplementedOrderServer
	chargeClient pb.ChargeClient
}

func (s *orderServer) PlaceOrder(ctx context.Context, req *pb.OrderRequest) (*pb.OrderResponse, error) {
	log.Printf("Order received: %d x %s", req.Quantity, req.Product)

	// Call the Charge Server
	chargeResponse, err := s.chargeClient.ChargeCustomer(ctx, &pb.ChargeRequest{
		CustomerId: "12345",                     // Hardcoded for simplicity
		Amount:     float32(req.Quantity) * 100, // Example calculation
	})

	if err != nil {
		appErr, err := errors.FromGRPCErr(err)
		if err != nil {
			log.Print("internal server error")
			return nil, err
		}
		log.Printf("appErr: %#v", appErr)
		if grpcErr := appErr.GetGRPCErr(); grpcErr != nil {
			log.Printf("Type of grpcErr: %s", reflect.TypeOf(grpcErr))
			switch grpcErr := grpcErr.(type) {
			case *proto.NotEnoughCharge:
				// Handle the specific error type *proto.NotEnoughCharge
				log.Println("Received NotEnoughCharge error")
				// Do something specific to this error
			default:
				// Handle other types or unexpected errors
				log.Printf("Received unexpected GRPC error: %T", grpcErr)
			}
		} else {
			log.Println("No GRPC error found in appErr")
		}
		return nil, err
	}

	return &pb.OrderResponse{
		Message: fmt.Sprintf("Order placed for %d x %s. %s", req.Quantity, req.Product, chargeResponse.Message),
	}, nil

}

func main() {
	serviceName := "order-service"
	debugMode := true
	logger.SetupDefaultLogger(slog.LevelDebug, true)
	// Connect to the Charge Server
	chargeConn, err := grpc.Dial("localhost:50052", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("Failed to connect to Charge Server: %v", err)
	}
	defer chargeConn.Close()

	chargeClient := pb.NewChargeClient(chargeConn)

	// Start the Order Server
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(interceptor.UnaryServerInterceptor(serviceName, debugMode)),
	)
	pb.RegisterOrderServer(grpcServer, &orderServer{chargeClient: chargeClient})

	log.Println("Order Server is running on port 50051...")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
