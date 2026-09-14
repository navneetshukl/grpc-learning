package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"

	pb "grpc-learning/stage-10-deployment-rollback/gen"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/status"
)

const port = 50051

// GreetService implements the GreetServiceServer interface.
// It demonstrates production-safe deployment practices:
// - Backward-compatible methods
// - Health checking
// - Graceful shutdown
type GreetService struct {
	pb.UnimplementedGreetServiceServer
	healthServer *health.Server
}

// Greet handles v1 requests (backward compatible).
func (s *GreetService) Greet(ctx context.Context, req *pb.GreetRequest) (*pb.GreetResponse, error) {
	name := req.GetName()
	if err := validateRequest(name); err != nil {
		return nil, err
	}
	log.Printf("Received Greet request for: %s", name)

	result := fmt.Sprintf("Hello, %s!", name)
	return &pb.GreetResponse{Result: result}, nil
}

// GreetV2 handles v2 requests (new version, safe addition).
// Old clients continue to call Greet; new clients use GreetV2.
func (s *GreetService) GreetV2(ctx context.Context, req *pb.GreetRequestV2) (*pb.GreetResponseV2, error) {
	name := req.GetName()
	if err := validateRequest(name); err != nil {
		return nil, err
	}

	language := req.GetLanguage()
	if language == "" {
		language = "en"
	}

	log.Printf("Received GreetV2 request for: %s (language: %s)", name, language)

	result := fmt.Sprintf("Hello, %s!", name)
	localized := getVersionedGreeting(name, language)

	return &pb.GreetResponseV2{
		Result:              result,
		LocalizedGreeting:   localized,
		Version:             "v2",
	}, nil
}

// getVersionedGreeting builds a localized greeting for v2 requests.
func getVersionedGreeting(name, language string) string {
	if language == "" {
		return fmt.Sprintf("Hello, %s!", name)
	}
	return fmt.Sprintf("Hello, %s! (localized for %s)", name, language)
}

// Check implements the custom health check in GreetService.
func (s *GreetService) Check(ctx context.Context, req *pb.HealthCheckRequest) (*pb.HealthCheckResponse, error) {
	service := req.GetService()
	log.Printf("Received custom health check for service: %s", service)

	return &pb.HealthCheckResponse{
		Status: "SERVING",
	}, nil
}

// GracefulStop stops the server gracefully, allowing active RPCs to finish.
func (s *GreetService) GracefulStop(grpcServer *grpc.Server, wg *sync.WaitGroup) {
	log.Printf("Graceful shutdown initiated...")

	// Stop accepting new connections
	// Existing RPCs continue until completion (connection draining)
	go func() {
		defer wg.Done()
		grpcServer.GracefulStop()
	}()
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

	// Step 3: Create service instances
	greetService := &GreetService{
		healthServer: health.NewServer(),
	}

	// Step 4: Register services
	pb.RegisterGreetServiceServer(grpcServer, greetService)
	healthpb.RegisterHealthServer(grpcServer, greetService.healthServer)

	// Step 5: Start accepting connections
	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	// Step 6: Wait for shutdown signal (Ctrl+C)
	// This is the key to graceful shutdown and connection draining
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	// Step 7: Graceful shutdown
	// - Stop accepting new connections
	// - Allow active RPCs to complete
	// - Then close everything
	var wg sync.WaitGroup
	wg.Add(1)
	greetService.GracefulStop(grpcServer, &wg)
	wg.Wait()

	log.Printf("Server stopped gracefully")
}

// validateRequest checks for invalid input and returns an appropriate gRPC error.
func validateRequest(name string) error {
	if name == "" {
		return status.Error(codes.InvalidArgument, "name is required")
	}
	return nil
}

