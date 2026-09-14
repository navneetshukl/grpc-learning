package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	pb "grpc-learning/stage-08-first-streaming-proto/gen"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	address = "localhost:50051"
)

func main() {
	// Step 1: Dial the server
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to dial server: %v", err)
	}
	defer conn.Close()
	log.Printf("Connected to server at %s", address)

	// Step 2: Create a client stub
	client := pb.NewGreetServiceClient(conn)

	// Step 3: Build the request for GreetManyTimes
	name := "Alice"
	count := int32(5)
	if len(os.Args) > 1 {
		name = os.Args[1]
	}
	if len(os.Args) > 2 {
		fmt.Sscanf(os.Args[2], "%d", &count)
	}

	req := &pb.GreetManyTimesRequest{
		Name:  name,
		Count: count,
	}

	// Step 4: Call GreetManyTimes (server-streaming RPC)
	// The server will send multiple responses, and we receive them one by one.
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	stream, err := client.GreetManyTimes(ctx, req)
	if err != nil {
		log.Fatalf("GreetManyTimes failed: %v", err)
	}

	// Step 5: Receive and print each greeting as it arrives
	fmt.Printf("Receiving %d greetings for %s:\n", count, name)
	for {
		resp, err := stream.Recv()
		if err != nil {
			// io.EOF means the server has finished sending all responses
			if err.Error() == "EOF" {
				break
			}
			log.Fatalf("Error receiving greeting: %v", err)
		}

		fmt.Println(resp.GetResult())
	}

	fmt.Println("Done receiving all greetings.")
}