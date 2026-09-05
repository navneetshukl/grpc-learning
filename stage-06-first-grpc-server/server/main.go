package main

import (
	"context"
	"fmt"
	"log"
	"net"

	pb "grpc-learning/stage-06-first-grpc-server/gen"

	"google.golang.org/grpc"
)

const port = 50051

// GreetService implements the GreetServiceServer interface.
// We embed UnimplementedGreetServiceServer so that if we forget to implement
// a method, the compiler will catch it at build time instead of crashing at runtime.
type GreetService struct {
	pb.UnimplementedGreetServiceServer
}

// Greet handles incoming Greet RPC requests.
func (s *GreetService) Greet(ctx context.Context, req *pb.GreetRequest) (*pb.GreetResponse, error) {
	name := req.GetName()
	log.Printf("Received Greet request for: %s", name)

	result := fmt.Sprintf("Hello, %s!", name)
	return &pb.GreetResponse{Result: result}, nil
}

func main() {
	// Step 1: Open a TCP port for listening
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("Failed to listen on port %d: %v", port, err)
	}
	log.Printf("Server listening on port %d", port)

	// Step 2: Create a new gRPC server
	grpcServer := grpc.NewServer()

	// Step 3: Register our GreetService with the gRPC server.
	// The generated code provides RegisterGreetServiceServer. It wires up
	// the handler so gRPC knows which method to call for each request.
	pb.RegisterGreetServiceServer(grpcServer, &GreetService{})

	// Step 4: Tell the gRPC server to start accepting connections.
	// This call blocks forever — the server runs until the process is killed.
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
