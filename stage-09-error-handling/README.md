# Stage 09 — Error Handling & gRPC Status Codes

> **Phase**: Hands-On Implementation  
> **Type**: Code stage — you will add comprehensive error handling  
> **Goal**: Learn how to properly handle errors in gRPC using status codes and demonstrate common error scenarios

---

## What I'm Learning

1. **gRPC error model** — using `status.Error()` with standardized codes
2. **Error classification** — what each code means and when to use it
3. **Client-side error handling** — inspecting status codes and messages
4. **Common gRPC errors** — NotFound, InvalidArgument, DeadlineExceeded, etc.
5. **Error propagation** — passing error details through streaming and unary RPCs

---

## Why Error Handling Matters in gRPC

Unlike REST APIs where you might return HTTP status codes, gRPC has its own error model built on top of HTTP/2:

- **Status codes** are standardized (`codes.NotFound`, `codes.InvalidArgument`, etc.)
- **Error details** can include custom messages and structured data
- **Streaming errors** require different handling — streams can fail mid-way
- **Client consistency** — both client and server agree on error types

---

## The gRPC Error Model

### Server-side: Throwing Errors

```go
// Return a standardized gRPC error with code and message
return nil, status.Error(codes.NotFound, "user not found")

// You can add additional details
return nil, status.Errorf(
    codes.InvalidArgument,
    "name must be at least 2 characters, got %d", len(name),
)
```

### Client-side: Catching Errors

```go
resp, err := client.GetError(ctx, req)
if err != nil {
    // Convert gRPC error back to status for inspection
    s := status.Convert(err)
    switch s.Code() {
    case codes.NotFound:
        fmt.Println("Resource was not found")
    case codes.InvalidArgument:
        fmt.Println("Client sent bad request")
    case codes.DeadlineExceeded:
        fmt.Println("Operation took too long")
    }
}
```

---

## Common gRPC Status Codes

| Code | Meaning | When to Use |
|------|---------|-------------|
| **OK** | 0 | Success (usually implicit) |
| **NotFound** | 5 | Resource doesn't exist (user, item, file) |
| **InvalidArgument** | 3 | Client sent bad input (empty name, wrong format) |
| **DeadlineExceeded** | 4 | Server took too long (timeout) |
| **Unimplemented** | 12 | Method not implemented (future feature) |
| **Unavailable** | 14 | Service temporarily down (maintenance, overload) |
| **Internal** | 13 | Unexpected server error |
| **PermissionDenied** | 7 | Client doesn't have permission |

---

## How to Run This Stage

### Step 1: Generate Go code

```bash
cd /Users/navneetshukla/Desktop/Projects/test/grpc/stage-09-error-handling
./generate.sh
```

### Step 2: Start the server

```bash
cd server
go run main.go
# Keep running in Terminal 1
```

### Step 3: Run the client

```bash
cd ../client
go run main.go
# Tests all error handling scenarios
```

---

## Key Code Changes

### Server-side error handling

```go
func (s *GreetService) GetError(ctx context.Context, req *pb.GreetErrorRequest) (*pb.GreetErrorResponse, error) {
    errorType := req.GetErrorType()
    
    switch errorType {
    case "not_found":
        return nil, status.Error(codes.NotFound, "user not found")
    
    case "invalid_argument":
        return nil, status.Error(codes.InvalidArgument, "invalid name provided")
    
    case "deadline_exceeded":
        time.Sleep(10 * time.Second) // Simulate timeout
        return nil, status.Error(codes.DeadlineExceeded, "operation timed out")
    
    // ... other error types
    }
}
```

### Client-side error inspection

```go
_, err := client.GetError(ctx, &pb.GreetErrorRequest{ErrorType: "not_found"})
if err != nil {
    s := status.Convert(err)
    if s.Code() == codes.NotFound {
        fmt.Println("✓ Correctly identified as NOT_FOUND")
    }
}
```

---

## What to Observe

- [ ] Server returns `codes.NotFound` for "not_found" error type
- [ ] Client correctly identifies and displays each error code
- [ ] Streaming still works when unary RPC fails
- [ ] Error messages provide actionable information

---

## Files in This Stage

```
stage-09-error-handling/
├── README.md          ← (you are here)
├── go.mod             ← Go module
├── generate.sh        ← Run protoc
├── proto/greet.proto  ← Contract with error handling RPC
├── gen/               ← Generated code (after protoc)
├── server/main.go     ← Server with comprehensive error handling
└── client/main.go     ← Client testing all error scenarios
```

---

## What is Next

**Stage 10 — Security: TLS & Authentication**

- Generate self-signed certificates
- Configure server TLS credentials
- Configure client certificate verification
- Add metadata/token-based authentication
- Implement mTLS for service-to-service communication

**After you complete this stage, say ready for stage 10**