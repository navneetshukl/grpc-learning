# Stage 08 — Server-Streaming RPCs

> **Phase**: Hands-On Implementation  
> **Type**: Code stage — you will add streaming RPC  
> **Goal**: Add a server-streaming method and see multiple responses arrive as a stream

---

## What I'm Learning

1. How **server-streaming** works: one request, many responses
2. The `stream` parameter in the generated server interface
3. How the client **receives** and iterates over responses
4. Distinguishing streaming types: unary, server, client, bidirectional

---

## Streaming Types Quick Reference

| Type | Pattern | Use Case |
|------|---------|----------|
| **Unary** | RPC(req) → resp | Simple request/response |
| **Server-Streaming** | RPC(req) → stream(resp) | Bulk data, logs, results |
| **Client-Streaming** | stream(req) → RPC(resp) | File upload, bulk send |
| **Bidi-Streaming** | stream(req) ↔ stream(resp) | Chat, real-time sync |

---

## How to Run This Stage

### Step 1: Generate Go code

```bash
cd /Users/navneetshukla/Desktop/Projects/test/grpc/stage-08-first-streaming-proto
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
go run main.go Alice 5
# Output will show 5 streaming greetings
```

---

## Key Code Changes

### Server-side (streaming handler)

```go
func (s *GreetService) GreetManyTimes(req *pb.GreetManyTimesRequest, 
    stream grpc.ServerStreamingServer[pb.GreetManyTimesResponse]) error {

    for i := 1; i <= count; i++ {
        if err := stream.Send(&pb.GreetManyTimesResponse{Result: result}); err != nil {
            return err
        }
    }
    return nil
}
```

### Client-side (receiving stream)

```go
stream, err := client.GreetManyTimes(ctx, req)
for {
    resp, err := stream.Recv()
    if err == io.EOF { break }  // Stream finished
    fmt.Println(resp.GetResult())
}
```

---

## What to Observe

- [ ] Server sends responses one at a time (1 second apart)
- [ ] Client receives each response via `stream.Recv()`
- [ ] `io.EOF` signals the end of the stream
- [ ] The stream is a Go channel-like construct over HTTP/2

---

## Files in This Stage

```
stage-08-first-streaming-proto/
├── README.md          ← (you are here)
├── go.mod             ← Go module
├── generate.sh        ← Run protoc
├── proto/
│   └── greet.proto    ← Contract with streaming method
├── gen/               ← Generated Go code (after protoc)
├── server/
│   └── main.go        ← Server with GreetManyTimes (streaming)
└── client/
    └── main.go        ← Client receiving the stream
```

---

## What is Next

**Stage 09 — Error Handling & Status Codes**

- Use `status.Error(codes.NotFound, "message")` on server
- Inspect error codes on client
- Add non-streaming error handling first