package main

import (
	"context"
	"fmt"
	"log"
	"net"

	pb "grpc-test/proto" // Replace with the correct import path

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
		return nil, fmt.Errorf("failed to charge customer: %v", err)
	}

	return &pb.OrderResponse{
		Message: fmt.Sprintf("Order placed for %d x %s. %s", req.Quantity, req.Product, chargeResponse.Message),
	}, nil
}

func main() {
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

	grpcServer := grpc.NewServer()
	pb.RegisterOrderServer(grpcServer, &orderServer{chargeClient: chargeClient})

	log.Println("Order Server is running on port 50051...")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
