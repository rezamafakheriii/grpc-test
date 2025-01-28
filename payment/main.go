package main

import (
	"context"
	"log"
	"log/slog"
	"math/rand"
	"net"
	"time"

	"grpc-test/domain"
	"grpc-test/interceptor"
	pb "grpc-test/proto" // Replace with the correct import path

	logger "github.com/revotech-group/go-lib/log"
	"google.golang.org/grpc"
)

type chargeServer struct {
	pb.UnimplementedChargeServer
}

func (s *chargeServer) ChargeCustomer(ctx context.Context, req *pb.ChargeRequest) (*pb.ChargeResponse, error) {
	log.Printf("Charge request received for customer %s: $%.2f", req.CustomerId, req.Amount)

	// Seed the random number generator
	rand.Seed(time.Now().UnixNano())

	// Randomly choose one of the two errors
	if rand.Intn(2) == 0 {
		return nil, domain.ErrNotEnoughCredit()
	} else {
		return nil, domain.ErrGatewayNotReachable()
	}
}

func main() {
	serviceName := "payment-service"
	debugMode := true
	logger.SetupDefaultLogger(slog.LevelDebug, true)
	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(interceptor.UnaryServerInterceptor(serviceName, debugMode)),
	)
	pb.RegisterChargeServer(grpcServer, &chargeServer{})

	log.Println("Charge Server is running on port 50052...")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
