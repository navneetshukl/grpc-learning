package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	pb "grpc-learning/stage-07-first-grpc-client/gen"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	address = "localhost:50051"
)

func main() {
	// Step 1: Dial the server — establish a connection
	// grpc.Dial does not block; it returns immediately. The actual connection
	// is established lazily, on the first RPC call.
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to dial server: %v", err)
	}
	defer conn.Close()
	log.Printf("Connected to server at %s", address)

	// Step 2: Create a client stub from the connection
	// NewGreetServiceClient takes the connection and returns a GreetServiceClient.
	client := pb.NewGreetServiceClient(conn)

	// Step 3: Build the request
	name := "Alice"
	if len(os.Args) > 1 {
		name = os.Args[1]
	}
	req := &pb.GreetRequest{Name: name}

	// Step 4: Call the Greet method
	// We set a timeout via context so the client doesn't hang forever.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := client.Greet(ctx, req)
	if err != nil {
		log.Fatalf("Greet failed: %v", err)
	}

	// Step 5: Print the response
	fmt.Printf("Server says: %s\n", resp.GetResult())
}
