package main

import (
	"log"
	"math/rand"
	"net"
	"time"

	pb "grpc-test/proto" // Replace with the correct import path

	"github.com/revotech-group/go-lib/grpc/interceptors"
	"google.golang.org/grpc"
)

type currencyServer struct {
	pb.UnimplementedCurrencyServer
}

func (s *currencyServer) SendExchangeRates(stream pb.Currency_SendExchangeRatesServer) error {
	currencies := []string{"USD", "EUR", "GBP", "JPY", "AUD"}

	for {
		// Simulate fetching exchange rates
		from := currencies[rand.Intn(len(currencies))]
		to := currencies[rand.Intn(len(currencies))]
		rate := rand.Float64() * (rand.Float64() + 0.5) // Random exchange rate for demo purposes

		exchangeRate := &pb.ExchangeRate{
			CurrencyFrom: from,
			CurrencyTo:   to,
			Rate:         rate,
			Timestamp:    time.Now().Format(time.RFC3339),
		}

		// Send the exchange rate to the payment service
		if err := stream.SendMsg(exchangeRate); err != nil {
			log.Printf("Error sending exchange rate: %v", err)
			return err
		}

		// Sleep before sending the next exchange rate
		time.Sleep(5 * time.Second) // Send exchange rates every 5 seconds
	}
}

func main() {
	// Setup the gRPC server
	lis, err := net.Listen("tcp", ":50053")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer(
		grpc.StreamInterceptor(interceptors.StreamServerErrorInterceptor()),
	)
	pb.RegisterCurrencyServer(grpcServer, &currencyServer{})

	log.Println("Currency Server is running on port 50055...")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
