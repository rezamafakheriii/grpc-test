package main

import (
	"context"
	"io"
	"log"
	"log/slog"
	"math/rand"
	"net"
	"time"

	"grpc-test/domain"
	"grpc-test/proto"
	pb "grpc-test/proto" // Replace with the correct import path

	"github.com/revotech-group/go-lib/grpc/interceptors"
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

func subscribeToExchangeRates() {
	// Establish connection to the currency service
	conn, err := grpc.Dial("localhost:50053", grpc.WithInsecure(), grpc.WithStreamInterceptor(interceptors.ClientStreamErrorInterceptor))
	if err != nil {
		log.Fatalf("Failed to connect to currency service: %v", err)
		return
	}
	defer conn.Close()

	// Create a new currency client
	client := pb.NewCurrencyClient(conn)

	// Subscribe to exchange rates from the currency service
	stream, err := client.SendExchangeRates(context.Background())
	if err != nil {
		log.Fatalf("Failed to subscribe to exchange rates: %v", err)
		return
	}

	// Continuously receive exchange rates
	for {
		var exchangeRate proto.ExchangeRate
		err := stream.RecvMsg(&exchangeRate)
		if err != nil {
			// Check for EOF or stream closure
			if err == io.EOF {
				log.Println("Stream closed by server.")
				break // Gracefully break out of the loop
			}
			log.Fatalf("Failed to receive exchange rate: %v", err)
			return
		}

		// Process the received exchange rate
		log.Printf("Received exchange rate: %s to %s = %.4f", exchangeRate.CurrencyFrom, exchangeRate.CurrencyTo, exchangeRate.Rate)
	}
}

func main() {
	// go subscribeToExchangeRates()

	logger.SetupDefaultLogger(slog.LevelDebug, true)
	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(interceptors.UnaryServerErrorInterceptor()),
	)
	pb.RegisterChargeServer(grpcServer, &chargeServer{})

	log.Println("Charge Server is running on port 50052...")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
