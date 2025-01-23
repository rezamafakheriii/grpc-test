package main

import (
	"context"
	"log"
	"time"

	pb "grpc-test/proto"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewProductOrderServiceClient(conn)

	// List products
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// Place an order
	order := &pb.OrderRequest{
		CustomerId: "customer-1",
		Items: []*pb.OrderItem{
			{ProductId: "1", Quantity: 2},
			{ProductId: "3", Quantity: 1},
		},
	}

	orderResponse, err := client.PlaceOrder(ctx, order)
	if err != nil {
		errorStatus := status.Convert(err)
		for _, d := range errorStatus.Details() {
			switch info := d.(type) {
			case *errdetails.BadRequest:
				log.Printf("Request Field Invalid: %s", info)
			default:
				log.Printf("Unexpected error type: %s", info)
			}
		}
		log.Print(errorStatus)
		return
	}

	log.Printf("Order placed! ID: %s, Total Price: %.2f", orderResponse.OrderId, orderResponse.TotalPrice)
}
