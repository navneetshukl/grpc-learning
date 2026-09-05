# Stage 03 — Client-Server Communication in gRPC

> **Phase**: Foundations (Read Only)
> **Type**: Concept-only stage — no code, no Go, no files to run
> **Goal**: Understand exactly what happens when a client calls a server method, step by step
> **Prerequisites**: Stage 01, Stage 02

---

## What I am Learning

1. The full lifecycle of a gRPC call from client to server and back
2. What a **client stub** actually is and why the developer does not touch the network
3. What a **gRPC channel** does
4. The role of HTTP/2 in transporting gRPC messages
5. What happens on the server side when a request arrives

---

## What is New vs Stage 02

Stage 01 introduced what gRPC is. Stage 02 compared it to REST. This stage zooms into the **mechanics** of one gRPC call. We trace every step from the moment you call a method until you get a response back.

No new files. No new code. Just understanding the request flow.

---

## The Four Actors in a gRPC Call

Every gRPC system has four actors. Each has a clear, separate job.

### 1. The Client Application (Your Code)

This is the Go program you write. It holds a **client stub** object. When you call a method on it, the network stuff happens automatically under the hood. To you, it looks exactly like calling a local function.

### 2. The Client Stub (Generated Code)

This is a Go struct generated from your `.proto` file by `protoc-gen-go-grpc`. It has one method per RPC you defined.

When you call `stub.Greet(ctx, req)`, the stub does these things automatically:

1. Marshals your request message to binary Protobuf
2. Builds an HTTP/2 request frame
3. Sends it over the channel
4. Waits for the response
5. Unmarshals the response back to a Go struct
6. Returns the result to your code

You never see any of this network code. The stub hides everything.

### 3. The gRPC Channel (Connection)

A channel is a long-lived connection to a server. When you call `grpc.Dial("localhost:50051", ...)`, you get a channel.

The channel:
- Manages the underlying TCP connection
- Speaks HTTP/2
- Handles reconnection if the server restarts
- Multiplexes many concurrent calls over one TCP connection

**Key point**: A channel is not a single call. It is a persistent connection. You create it once at startup and reuse it for all calls throughout your program's lifetime.

### 4. The Server (Generated + Your Code)

The server runs a **gRPC server instance** that listens on a port. When a request arrives:

1. The server reads the HTTP/2 frame from the network
2. It unmarshals the binary Protobuf into a Go struct
3. It looks up which method was called and calls your handler function
4. Your handler runs and returns a response
5. The server marshals the response to binary Protobuf
6. It sends the response back over HTTP/2

---

## The Complete Request-Response Flow

Here is the full step-by-step journey of a unary gRPC call. Read this carefully — it is the most important diagram in this stage.

```
Step 1: Client creates a channel
┌─────────────────────────────────────────────────────┐
│  conn, err := grpc.Dial("localhost:50051", opts)   │
└──────────────────┬──────────────────────────────────┘
                   │ TCP + HTTP/2 connection established
                   ▼
         ┌─────────────────┐
         │  gRPC Channel   │  ← persists for the lifetime of the program
         │  (localhost:50051)│
         └────────┬────────┘

Step 2: Client creates a stub from the channel
┌─────────────────────────────────────────────────────┐
│  client := greetpb.NewGreetServiceClient(conn)      │
└──────────────────┬──────────────────────────────────┘
                   │ stub holds a reference to the channel
                   ▼

Step 3: Client calls a method on the stub
┌─────────────────────────────────────────────────────┐
│  req := &greetpb.GreetRequest{Name: "Alice"}       │
│  resp, err := client.Greet(ctx, req)               │
└──────────────────┬──────────────────────────────────┘
                   │ Stub marshals: GreetRequest → binary Protobuf
                   ▼

Step 4: Stub sends over HTTP/2
┌─────────────────────────────────────────────────────┐
│  HTTP/2 POST frame                                 │
│  Path: /greet.GreetService/Greet                   │
│  Body: <binary Protobuf data>                      │
│  Headers: Content-Type: application/grpc           │
└──────────────────┬──────────────────────────────────┘
                   │ TCP/IP packets travel over the network
                   ▼

Step 5: Server receives the frame
┌─────────────────────────────────────────────────────┐
│  gRPC Server listening on :50051                    │
│  Reads HTTP/2 frame, unmarshals binary → Go struct   │
└──────────────────┬──────────────────────────────────┘
                   │ Server looks up "Greet" method
                   ▼

Step 6: Server calls your handler
┌─────────────────────────────────────────────────────┐
│  func (s *server) Greet(ctx, req) (*resp, error) {  │
│      name := req.GetName()  // "Alice"              │
│      result := "Hello, " + name                     │
│      return &greetpb.GreetResponse{Result: result}, nil
│  }                                                  │
└──────────────────┬──────────────────────────────────┘
                   │ Response marshalled to binary Protobuf
                   ▼

Step 7: Response sent back over HTTP/2
┌─────────────────────────────────────────────────────┐
│  HTTP/2 frame with binary Protobuf body             │
│  Status: 200 OK                                    │
└──────────────────┬──────────────────────────────────┘
                   │ TCP/IP back to client
                   ▼

Step 8: Client stub unmarshals and returns
┌─────────────────────────────────────────────────────┐
│  resp.GetResult()  // "Hello, Alice"               │
└─────────────────────────────────────────────────────┘
```

The entire journey from step 3 to step 8 happens in milliseconds. Your code at step 3 waits until step 8 is complete. To you, it feels exactly like a local function call.


---

## A Simpler View: The Abstraction Layers

Think of a gRPC call as four layers stacked on top of each other. Each layer only talks to the layer directly above or below it.

```
┌──────────────────────────────────────┐
│         Your Application Code        │  ← you write this
│  client.Greet(ctx, req)              │
└──────────────────┬───────────────────┘
                   │  looks like a function call
┌──────────────────▼───────────────────┐
│         Generated Client Stub        │  ← generated by protoc
│  - Marshals request to Protobuf      │
│  - Unmarshals response from Protobuf │
└──────────────────┬───────────────────┘
                   │  sends binary data
┌──────────────────▼───────────────────┐
│         gRPC Library (Go)            │  ← grpc-go package
│  - Manages HTTP/2 frames             │
│  - Handles multiplexing              │
│  - Manages the TCP connection        │
└──────────────────┬───────────────────┘
                   │  raw bytes over TCP
┌──────────────────▼───────────────────┐
│         TCP/IP Network               │  ← the actual wire
│  - Sends bytes from client to server │
└──────────────────────────────────────┘
```

The key insight: **your code never touches the network directly**. It talks to the stub. The stub talks to the gRPC library. The gRPC library talks to TCP. Each layer has one job.

---

## The Channel is a Long-Lived Connection

One of the most important things to understand: the channel is not created per request. It is created once and reused.

**WRONG way (do not do this):**
```go
for _, name := range names {
    conn, _ := grpc.Dial("localhost:50051")   // new connection every call!
    client := greetpb.NewGreetServiceClient(conn)
    client.Greet(ctx, &req)
    conn.Close()
}
```

**RIGHT way (do this):**
```go
conn, _ := grpc.Dial("localhost:50051")        // one connection
client := greetpb.NewGreetServiceClient(conn)
for _, name := range names {
    client.Greet(ctx, &req)                    // reuse the connection
}
conn.Close()  // close only when shutting down
```

HTTP/2's multiplexing means one TCP connection can carry thousands of concurrent gRPC calls. Opening a new connection for every call is wasteful and slow.

---

## Where Protobuf Fits In

Notice in the flow above that the request is marshalled to binary Protobuf before being sent. Protobuf is the **encoding format** — the language that both sides use to represent data on the wire.

```
Your Go struct:  GreetRequest{Name: "Alice"}
                         │
                         │  marshalled by the stub
                         ▼
Wire format:  <binary Protobuf bytes>
                         │
                         │  unmarshalled by the server
                         ▼
Server Go struct:  GreetRequest{Name: "Alice"}
```

The binary format is compact and fast to parse. Field names are not sent — only field numbers. This is why Protobuf messages are much smaller than JSON.

---

## Files in This Stage

**None.** This is a reading stage.

```
stage-03-client-server-model/
└── README.md   (you are here)
```

---

## How to Run This Stage

There is nothing to execute. Instead:

1. **Read** this README carefully, especially the 8-step flow diagram
2. **Trace through** the flow with a mental example: what if the client called `Greet` with `Name: "Bob"`?
3. **Memorize** the four layers of the abstraction stack

---

## What to Observe / Reflect On

After reading, you should be able to answer in your own words:

- [ ] What are the four actors in a gRPC call?
- [ ] What does the client stub do when you call a method?
- [ ] Why is the channel created once and reused?
- [ ] What would break if the client stub did not exist and you had to send data over the network yourself?
- [ ] In step 5 of the flow, how does the server know which method to call?
- [ ] Why does the channel use HTTP/2 instead of raw TCP?

---

## Exercise

There is no code exercise for this stage. Instead:

1. **Draw the 8-step flow from memory** — start from `grpc.Dial`, then channel, stub, method call, marshalling, HTTP/2 send, server receive, handler call, response marshalling, and back to the client.

2. **Draw the 4-layer abstraction stack** — application code, stub, gRPC library, TCP/IP. Label what each layer is responsible for.

3. **Explain to someone** (or write down): what happens when you call `client.Greet(ctx, req)`? Start from that line of code and describe the full journey until the response comes back.

---

## What is Next

**Stage 04 — Writing Your First .proto File**

We finally start writing code. You will:
- Write a `.proto` file by hand
- Define messages and services
- Understand field numbers and field types
- Prepare to generate Go code from the proto definition

After you complete this stage, say **"ready for stage 04"** and we will continue.
