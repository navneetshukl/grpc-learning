# Stage 10 — Deployment Strategies & Rollback in gRPC

> **Phase**: Production Readiness (Hands-On)
> **Type**: Code stage — you will implement deployment-safe patterns
> **Goal**: Learn how to deploy and rollback gRPC services safely using versioning, health checks, and graceful shutdown
> **Prerequisite**: Stage 09

---

## What I am Learning

In this stage, you will:

1. How to version gRPC services safely using Protobuf backward compatibility
2. How to run multiple versions of a service on the same port
3. How to implement graceful shutdown with connection draining
4. How to use the gRPC health checking protocol for service discovery
5. How to test rollback scenarios by switching between service versions
6. How to validate that old clients continue to work during updates

---

## What is New vs Stage 09

Stage 09 focused on error handling — returning proper status codes when things go wrong. This stage focuses on **deployment safety** — how to update and rollback services without breaking clients.

We'll implement:
- **Versioned services**: v1 and v2 methods on the same port
- **Backward-compatible Protobuf**: Adding fields without breaking existing clients
- **Health checking**: Using the standard gRPC health protocol
- **Graceful shutdown**: Allowing active RPCs to finish before stopping
- **Connection draining**: Stopping new connections while completing existing ones

---

## Why This Exists

In production, services are updated frequently. The challenge is:
- How do you update a service without breaking existing clients?
- How do you rollback if something goes wrong?
- How do you know if your service is healthy during deployment?

gRPC provides built-in tools to solve these problems:
- **Protobuf evolution rules** let you add fields safely
- **Multiple service versions** can coexist on the same port
- **Health checks** let load balancers know when it's safe to send traffic
- **Graceful shutdown** prevents dropping active connections

This stage shows you how to use these features together for safe deployments.
---

## How to Run This Stage

### Step 1: Generate Go code

```bash
cd /Users/navneetshukla/Desktop/Projects/test/grpc/stage-10-deployment-rollback
./generate.sh
```

This compiles `proto/greet.proto` into Go code in the `gen/` directory.

### Step 2: Start the server

```bash
cd server
go run main.go
# Keep this running in Terminal 1
```

You should see:
```
2026/09/14 21:55:00 Server listening on port 50051
```

### Step 3: Run the client

In a second terminal:

```bash
cd ../client
go run main.go
```

You should see output showing:
- Successful v1 Greet RPC calls (backward compatible)
- Successful v2 GreetV2 RPC calls (enriched response)
- Health check passing (SERVING status)
- Version-aware deployment summary

### Step 4: Test graceful shutdown

In the server terminal, press `Ctrl+C` (or send SIGTERM):
```
2026/09/14 21:55:30 Received shutdown signal, initiating graceful stop...
2026/09/14 21:55:30 Graceful shutdown initiated...
2026/09/14 21:55:30 Server stopped gracefully
```

The server will stop accepting new connections but allow existing RPCs to complete.

---

## What Each Part Does

### Protobuf Versioning Strategy (`proto/greet.proto`)

The `.proto` file demonstrates safe evolution:
- **v1 messages**: `GreetRequest` and `GreetResponse` (original contract)
- **v2 messages**: `GreetRequestV2` and `GreetResponseV2` (add new fields safely)
- **Field numbers**: Never change existing field numbers; new fields get higher numbers
- **Backward compatibility**: v1 clients ignore v2 fields; v2 clients can still use v1 methods

### Server Implementation (`server/main.go`)

1. **Service struct**: Embeds `UnimplementedGreetServiceServer` for forward compatibility
2. **Health server**: Uses gRPC's built-in health checking protocol
3. **Greet method**: Handles v1 requests (original contract)
4. **GreetV2 method**: Handles v2 requests (enriched with localization and version)
5. **GracefulStop**: Implements connection draining via `grpcServer.GracefulStop()`
6. **Signal handling**: Waits for SIGINT/SIGTERM to initiate graceful shutdown
7. **Validation**: Returns proper gRPC error codes for invalid input

### Client Implementation (`client/main.go`)

1. **Connection**: Dials the server with insecure credentials (for learning)
2. **v1 testing**: Calls the original `Greet` method to verify backward compatibility
3. **v2 testing**: Calls the enhanced `GreetV2` method with language parameter
4. **Health checking**: Uses the standard gRPC health protocol to check service status
5. **Version summary**: Prints deployment readiness verification

### Generated Code (`gen/`)

Contains the compiled Protobuf output:
- `greet.pb.go`: Message structs (GreetRequest, GreetResponse, etc.)
- `greet_grpc.pb.go`: Client stubs and server interfaces

### Build Files

- `go.mod`: Defines the Go module and dependencies
- `go.sum`: Dependency checksums for reproducible builds
- `generate.sh`: Script to regenerate Go code from `.proto` files
---

## Files in This Stage

```
stage-10-deployment-rollback/
├── README.md          ← (you are here)
├── go.mod             ← Go module file
├── go.sum             ← Go module checksums
├── proto/
│   └── greet.proto    ← Versioned contract showing v1/v2 methods
├── gen/
│   ├── greet.pb.go        ← Generated message types
│   └── greet_grpc.pb.go   ← Generated client/server stubs
├── server/
│   └── main.go    ← Server with versioning, health checks, graceful shutdown
├── client/
│   └── main.go    ← Client testing v1/v2 compatibility and health checks
└── generate.sh      ← Runs protoc to generate Go code
```

---

## What to Observe / Reflect On

After completing this stage, you should be able to explain:

- [ ] How Protobuf field numbering enables backward compatibility when adding new fields
- [ ] Why we can have both v1 and v2 methods on the same service without conflict
- [ ] How the gRPC health checking protocol helps with service discovery and load balancing
- [ ] What connection draining means and why it's important during deployments
- [ ] How `grpcServer.GracefulStop()` differs from immediately calling `grpcServer.Stop()`
- [ ] Why embedding `UnimplementedGreetServiceServer` provides forward compatibility
- [ ] How versioned services enable blue/green or rolling deployment strategies
- [ ] What happens to in-flight RPCs during a graceful shutdown
- [ ] How health checks prevent sending traffic to unhealthy instances during deployment

---

## Exercise

1. **Add a v3 version**: Create `GreetRequestV3` and `GreetResponseV3` with additional fields, and implement a `GreetV3` method on the server. Test that v1 and v2 clients still work.

2. **Simulate a rollback**: 
   - Start the server with v2 only (comment out v1 methods in the proto)
   - Run the client to verify it works with v2
   - "Rollback" by restoring v1 methods and removing v2
   - Restart the server and verify v1 clients still work
   - (Hint: You'll need to regenerate the code after each change)

3. **Add latency to health checks**: Modify the health check implementation to randomly return `NOT_SERVING` 10% of the time. Observe how the client handles this.

4. **Test connection draining**: 
   - Modify the client to make long-running requests (add a sleep in the request)
   - Start multiple concurrent client connections
   - Initiate graceful shutdown and observe that the server waits for active RPCs to complete

5. **Implement TLS**: Add TLS credentials to both client and server to secure the connection in production-like conditions.

---

## What is Next

Congratulations! You've completed all 10 stages of the gRPC learning path. You now understand:

- **Foundations**: What gRPC is and why it exists (Stages 01-03)
- **Contract-first development**: Writing `.proto` files and generating code (Stages 04-05)
- **Implementation**: Building unary and streaming servers and clients (Stages 06-08)
- **Production readiness**: Error handling, deployment strategies, and observability (Stages 09-10)

From here, you can:
- Explore advanced gRPC features like reflection, interceptors, or custom metadata
- Integrate gRPC with service meshes like Istio or Linkerd
- Experiment with other languages (Java, Python, Node.js) using the same `.proto` contracts
- Build production microservices with gRPC as the communication layer

You have a solid foundation in gRPC — go build something amazing! 🚀