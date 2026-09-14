package main

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"

	pb "grpc-learning/stage-09-error-handling/gen"
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

	// Step 3: Test the unary Greet RPC
	fmt.Println("=== Testing Greet RPC ===")
	req := &pb.GreetRequest{Name: "Alice"}
	resp, err := client.Greet(context.Background(), req)
	if err != nil {
		log.Fatalf("Greet failed: %v", err)
	}
	fmt.Printf("Greet response: %s\n", resp.GetResult())

	// Step 4: Test error handling RPCs
	fmt.Println("\n=== Testing Error Handling ===")

	// Test NOT_FOUND error
	fmt.Println("\n--- Testing NOT_FOUND error ---")
	_, err = client.GetError(context.Background(), &pb.GreetErrorRequest{ErrorType: "not_found"})
	if err != nil {
		// This is expected - we expect an error
		s := status.Convert(err)
		fmt.Printf("Error received (expected): Code=%v, Message=%s\n", s.Code(), s.Message())
		if s.Code() == codes.NotFound {
			fmt.Println("✓ Correctly identified as NOT_FOUND")
		}
	}

	// Test INVALID_ARGUMENT error
	fmt.Println("\n--- Testing INVALID_ARGUMENT error ---")
	_, err = client.GetError(context.Background(), &pb.GreetErrorRequest{ErrorType: "invalid_argument"})
	if err != nil {
		s := status.Convert(err)
		fmt.Printf("Error received (expected): Code=%v, Message=%s\n", s.Code(), s.Message())
		if s.Code() == codes.InvalidArgument {
			fmt.Println("✓ Correctly identified as INVALID_ARGUMENT")
		}
	}

	// Test DEADLINE_EXCEEDED error
	fmt.Println("\n--- Testing DEADLINE_EXCEEDED error ---")
	_, err = client.GetError(context.Background(), &pb.GreetErrorRequest{ErrorType: "deadline_exceeded"})
	if err != nil {
		s := status.Convert(err)
		fmt.Printf("Error received (expected): Code=%v, Message=%s\n", s.Code(), s.Message())
		if s.Code() == codes.DeadlineExceeded {
			fmt.Println("✓ Correctly identified as DEADLINE_EXCEEDED")
		}
	}

	// Test UNAVAILABLE error
	fmt.Println("\n--- Testing UNAVAILABLE error ---")
	_, err = client.GetError(context.Background(), &pb.GreetErrorRequest{ErrorType: "unavailable"})
	if err != nil {
		s := status.Convert(err)
		fmt.Printf("Error received (expected): Code=%v, Message=%s\n", s.Code(), s.Message())
		if s.Code() == codes.Unavailable {
			fmt.Println("✓ Correctly identified as UNAVAILABLE")
		}
	}

	// Test INTERNAL error
	fmt.Println("\n--- Testing INTERNAL error ---")
	_, err = client.GetError(context.Background(), &pb.GreetErrorRequest{ErrorType: "internal"})
	if err != nil {
		s := status.Convert(err)
		fmt.Printf("Error received (expected): Code=%v, Message=%s\n", s.Code(), s.Message())
		if s.Code() == codes.Internal {
			fmt.Println("✓ Correctly identified as INTERNAL")
		}
	}

	// Test the default case (no error)
	fmt.Println("\n--- Testing Default (no error) ---")
	defaultResp, err := client.GetError(context.Background(), &pb.GreetErrorRequest{ErrorType: "default"})
	if err != nil {
		log.Fatalf("GetError default failed: %v", err)
	}
	fmt.Printf("Response: Message=%s, Details=%s\n", defaultResp.GetMessage(), defaultResp.GetDetails())

	fmt.Println("\n=== All error handling tests complete ===")
}
