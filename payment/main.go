package main

import (
	"context"
	"fmt"
	"log"
	"net"

	pb "grpc-test/proto" // Replace with the correct import path

	"google.golang.org/grpc"
)

type chargeServer struct {
	pb.UnimplementedChargeServer
}

func (s *chargeServer) ChargeCustomer(ctx context.Context, req *pb.ChargeRequest) (*pb.ChargeResponse, error) {
	log.Printf("Charge request received for customer %s: $%.2f", req.CustomerId, req.Amount)
	return &pb.ChargeResponse{Message: fmt.Sprintf("Charged $%.2f to customer %s", req.Amount, req.CustomerId)}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterChargeServer(grpcServer, &chargeServer{})

	log.Println("Charge Server is running on port 50052...")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
