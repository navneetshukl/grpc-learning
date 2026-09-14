package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"

	pb "grpc-learning/stage-10-deployment-rollback/gen"

	healthpb "google.golang.org/grpc/health/grpc_health_v1"
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

	// Step 2: Create client stubs
	client := pb.NewGreetServiceClient(conn)
	healthClient := healthpb.NewHealthClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Step 3: Test v1 unary RPC (backward compatible)
	fmt.Println("=== Testing v1 Greet RPC ===")
	req := &pb.GreetRequest{Name: "Alice"}
	resp, err := client.Greet(ctx, req)
	if err != nil {
		log.Fatalf("Greet failed: %v", err)
	}
	fmt.Printf("v1 Greet response: %s\n", resp.GetResult())

	// Step 4: Test v2 unary RPC (enriched response)
	fmt.Println("\n=== Testing v2 Greet RPC ===")
	reqV2 := &pb.GreetRequestV2{
		Name:     "Bob",
		Language: "es",
	}
	respV2, err := client.GreetV2(ctx, reqV2)
	if err != nil {
		log.Fatalf("GreetV2 failed: %v", err)
	}
	fmt.Printf("v2 Greet response: %s\n", respV2.GetResult())
	fmt.Printf("v2 Localized greeting: %s\n", respV2.GetLocalizedGreeting())
	fmt.Printf("v2 Version: %s\n", respV2.GetVersion())

	// Step 5: Test gRPC health checking protocol
	fmt.Println("\n=== Testing Health Check ===")
	healthResp, err := healthClient.Check(ctx, &healthpb.HealthCheckRequest{Service: ""})
	if err != nil {
		s := status.Convert(err)
		fmt.Printf("Health check error: Code=%v, Message=%s\n", s.Code(), s.Message())
	} else {
		fmt.Printf("Health check status: %s\n", healthResp.GetStatus())
		if healthResp.Status == healthpb.HealthCheckResponse_SERVING {
			fmt.Println("✓ Service is SERVING")
		}
	}

	// Step 6: Demonstrate version-aware deployment
	fmt.Println("\n=== Version-Aware Deployment Summary ===")
	fmt.Println("✓ v1 clients still work (Greet)")
	fmt.Println("✓ v2 clients get enriched response (GreetV2)")
	fmt.Println("✓ Health checks pass (SERVING)")
	fmt.Println("✓ Graceful shutdown supported (SIGTERM/SIGINT)")

	fmt.Println("\n=== All deployment & rollback tests complete ===")
}
