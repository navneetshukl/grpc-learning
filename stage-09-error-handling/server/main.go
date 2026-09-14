package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	pb "grpc-learning/stage-09-error-handling/gen"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const port = 50051

// GreetService implements the GreetServiceServer interface.
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
func (s *GreetService) GreetManyTimes(req *pb.GreetManyTimesRequest,
	stream grpc.ServerStreamingServer[pb.GreetManyTimesResponse]) error {
	name := req.GetName()
	count := req.GetCount()

	log.Printf("Received GreetManyTimes request for: %s, count: %d", name, count)

	for i := int32(1); i <= count; i++ {
		result := fmt.Sprintf("Hello, %s! Greeting %d out of %d", name, i, count)

		if err := stream.Send(&pb.GreetManyTimesResponse{Result: result}); err != nil {
			log.Printf("Failed to send greeting %d: %v", i, err)
			return err
		}

		time.Sleep(1 * time.Second)
	}

	return nil
}

// GetError demonstrates error handling by returning specific gRPC status codes.
// This is useful for testing how clients handle different error scenarios.
func (s *GreetService) GetError(ctx context.Context, req *pb.GreetErrorRequest) (*pb.GreetErrorResponse, error) {
	errorType := req.GetErrorType()
	log.Printf("Received GetError request for: %s", errorType)

	switch errorType {
	case "not_found":
		// NOT_FOUND: The requested resource does not exist
		return nil, status.Error(codes.NotFound, "user not found")

	case "invalid_argument":
		// INVALID_ARGUMENT: Client provided invalid input
		return nil, status.Error(codes.InvalidArgument, "invalid name provided")

	case "deadline_exceeded":
		// DEADLINE_EXCEEDED: Server took too long to respond
		// Simulate a slow operation that exceeds the client's timeout
		time.Sleep(10 * time.Second)
		return nil, status.Error(codes.DeadlineExceeded, "operation timed out")

	case "internal":
		// INTERNAL: Server-side unexpected error
		return nil, status.Error(codes.Internal, "something went wrong on the server")

	case "unavailable":
		// UNAVAILABLE: Service temporarily unavailable
		return nil, status.Error(codes.Unavailable, "service is currently unavailable")

	default:
		return &pb.GreetErrorResponse{
			Message: "No error requested",
			Details: "Please specify error_type: not_found, invalid_argument, deadline_exceeded, internal, or unavailable",
		}, nil
	}
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
