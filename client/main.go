package main

import (
	"context"
	"log"
	"time"

	pb "grpc-test/proto" // Replace with the correct import path

	"google.golang.org/grpc"
)

func main() {
	// Connect to the Order Server
	orderConn, err := grpc.Dial("localhost:50051", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("Failed to connect to Order Server: %v", err)
	}
	defer orderConn.Close()

	orderClient := pb.NewOrderClient(orderConn)

	// Call the Order Server
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	orderResponse, err := orderClient.PlaceOrder(ctx, &pb.OrderRequest{
		Product:  "Laptop",
		Quantity: 2,
	})
	if err != nil {
		log.Fatalf("Failed to place order: %v", err)
	}
	log.Printf("Order Response: %s", orderResponse.Message)
}
