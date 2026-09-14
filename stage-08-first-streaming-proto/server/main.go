package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	pb "grpc-learning/stage-08-first-streaming-proto/gen"

	"google.golang.org/grpc"
)

const port = 50051

// GreetService implements the GreetServiceServer interface.
// We embed UnimplementedGreetServiceServer so that if we forget to implement
// a method, the compiler will catch it at build time instead of crashing at runtime.
type GreetService struct {
	pb.UnimplementedGreetServiceServer
}

// Greet handles incoming Greet RPC requests (unary RPC).
func (s *GreetService) Greet(ctx context.Context, req *pb.GreetRequest) (*pb.GreetResponse, error) {
	name := req.GetName()
	log.Printf("Received Greet request for: %s", name)

	result := fmt.Sprintf("Hello, %s!", name)
	return &pb.GreetResponse{Result: result}, nil
}

// GreetManyTimes handles server-streaming RPC requests.
// The server sends back multiple GreetManyTimesResponse messages for a single request.
// The method returns an error because the stream itself can fail (e.g., client disconnect).
func (s *GreetService) GreetManyTimes(req *pb.GreetManyTimesRequest, 
    stream grpc.ServerStreamingServer[pb.GreetManyTimesResponse]) error {
	name := req.GetName()
	count := req.GetCount()

	log.Printf("Received GreetManyTimes request for: %s, count: %d", name, count)

	// Loop to send multiple responses back to the client
	for i := int32(1); i <= count; i++ {
		result := fmt.Sprintf("Hello, %s! Greeting %d out of %d", name, i, count)

		// Send one response over the stream
		// If the client disconnects, this returns an error
		if err := stream.Send(&pb.GreetManyTimesResponse{Result: result}); err != nil {
			log.Printf("Failed to send greeting %d: %v", i, err)
			return err
		}

		// Small delay to make the streaming visible
		time.Sleep(1 * time.Second)
	}

	return nil
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
	pb.RegisterGreetServiceServer(grpcServer, &GreetService{})

	// Step 4: Tell the gRPC server to start accepting connections.
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}