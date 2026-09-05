# Stage 02 — gRPC vs REST

> **Phase**: Foundations (Read Only)
> **Type**: Concept-only stage — no code, no Go, no files to run
> **Goal**: Understand the deep differences between gRPC and REST so you know when to use which
> **Prerequisite**: Stage 01

---

## What I am Learning

1. What REST actually means (it is more than just JSON over HTTP)
2. A side-by-side comparison of gRPC and REST across every dimension
3. The specific tradeoffs of each approach
4. When to choose gRPC and when to choose REST
5. Why many real systems use BOTH

---

## What is New vs Stage 01

In Stage 01, we introduced gRPC as a concept. In this stage, we **compare it directly to REST** — the dominant API style you are likely already familiar with.

No new files. No new code. Just understanding the difference.

---

## Why This Exists

Most developers already know REST. So why learn gRPC?

The answer is: **they solve different problems**. REST is great for public-facing web APIs. gRPC is better for internal service-to-service communication, high-throughput systems, polyglot environments, and real-time streaming. Understanding both lets you make the right architectural choice.

---

## The Core Difference: How Requests Are Modeled

### REST: Resources and Verbs

REST models everything as **resources** with **verbs**:

```
GET    /users        → get all users
POST   /users        → create a user
GET    /users/42     → get user 42
PUT    /users/42     → update user 42
DELETE /users/42     → delete user 42
```

You interact with **nouns** (resources) using **HTTP verbs**. The URL identifies what. The verb identifies the action.

### gRPC: Services and Methods

gRPC models everything as **service methods**:

```
UserService.GetUser({id: 42})    → get user 42
UserService.CreateUser({...})     → create a user
UserService.UpdateUser({...})     → update a user
UserService.DeleteUser({id: 42})  → delete user 42
---

## Side-by-Side Comparison

### 1. Data Format

| | REST | gRPC |
|---|---|---|
| Format | JSON (text) | Protobuf (binary) |
| Human readable | Yes | No |
| Parsing speed | Slow (text) | Fast (binary) |
| Payload size | Large | Small |

gRPC payloads are roughly 3 to 10 times smaller than JSON and parse roughly 10 to 100 times faster.

### 2. Transport Protocol

| | REST | gRPC |
|---|---|---|
| Protocol | HTTP/1.1 or HTTP/2 | HTTP/2 only |
| Multiplexing | No (in HTTP/1.1) | Yes |
| Header compression | No | Yes (HPACK) |
| Bidirectional | No | Yes |

gRPC reuses a single HTTP/2 connection for many concurrent calls. REST with HTTP/1.1 opens a new connection per request.

### 3. API Contract

| | REST | gRPC |
|---|---|---|
| Contract | Optional (OpenAPI/Swagger) | Required (.proto file) |
| Enforced at build time | No | Yes |
| Generated client code | Optional extra step | Built-in |
| Schema evolution | Manual | Field numbers for compatibility |

With REST, you can ship a broken API and only find out at runtime. With gRPC, the compiler catches mismatches.

### 4. Browser Support

| | REST | gRPC |
|---|---|---|
| Direct browser calls | Yes | No |
| Requires gateway | No | Yes (gRPC-Web or gRPC-Gateway) |

If your API is consumed by web browsers, REST wins. gRPC needs an intermediary.

### 5. Streaming

| | REST | gRPC |
|---|---|---|
| Server streaming | No | Yes |
| Client streaming | No | Yes |
| Bidirectional streaming | No | Yes |

For real-time data, notifications, or large file transfers, REST requires workarounds like WebSockets or Server-Sent Events. gRPC has streaming built into the protocol.

### 6. Code Generation

| | REST | gRPC |
|---|---|---|
| Server stubs | Manual | Automatic from .proto |
| Client SDKs | Manual or optional | Automatic from .proto |
| Languages supported | Various generators | 11+ official, all identical |

With gRPC, one .proto file generates matching server plus client code in Go, Python, Java, Node, C#, and more, all at once.

### 7. Error Handling

| | REST | gRPC |
|---|---|---|
| Error codes | HTTP status codes (200, 404, 500) | gRPC status codes (NOT_FOUND, INVALID_ARGUMENT) |
| Error details | Freeform JSON message | Structured with code, message, details |

gRPC errors are richer and typed. REST errors depend on HTTP codes which were not designed for API logic.
```

---

## Real API Comparison

### REST Version

**Request:**
```http
GET /api/v1/users/42
Accept: application/json
Authorization: Bearer token123
```

**Response:**
```json
{
  "id": 42,
  "name": "Alice",
  "email": "alice@example.com",
  "age": 30
}
```

**Error (user not found):**
```
HTTP/1.1 404 Not Found
{
  "error": "User not found",
  "code": "USER_NOT_FOUND"
}
```

### gRPC Version

**Request:**
```go
req := &userpb.GetUserRequest{Id: 42}
ctx := metadata.NewOutgoingContext(ctx, md)
resp, err := client.GetUser(ctx, req)
```

**Success response:**
```go
resp := &userpb.GetUserResponse{
    User: &userpb.User{
        Id:    42,
        Name:  "Alice",
        Email: "alice@example.com",
        Age:   30,
    },
}
```

**Error (user not found):**
```go
// Server returns:
return nil, status.Errorf(codes.NotFound, "user 42 not found")

// Client receives:
err.Error()  // "rpc error: code = NotFound desc = user 42 not found"
st := status.FromError(err)
st.Code()    // codes.NotFound
st.Message() // "user 42 not found"
```
---

## Why REST Is Still Dominant

Given gRPC advantages, why does REST still dominate?

1. **Browser support** — Browsers cannot make gRPC calls directly. Most web APIs are consumed by browsers.
2. **Ecosystem** — REST has decades of tooling: Postman, curl, OpenAPI, Swagger, etc.
3. **Simplicity** — JSON is human-readable. You can debug it in the browser console.
4. **Caching** — HTTP caching works naturally with REST. gRPC caching is more complex.
5. **Industry momentum** — Most developers already know REST. gRPC has a steeper learning curve.

---

## When to Use What

Use **gRPC** when you need backend services talking to each other, high-throughput and low-latency, real-time streaming, polyglot environments, or internal APIs with strict contracts.

Use **REST** when you have public APIs consumed by web browsers, need caching, have simpler use cases with human-readable payloads, or when ecosystem and tooling matter more than performance.

Many real systems use **both** in the same architecture: external clients use REST (because browsers require it), while internal services use gRPC (because it is faster and more efficient). gRPC-Gateway bridges the gap.

---

## Real-World Architecture Pattern

In production, a common pattern is:

External clients (web, mobile, third-party) send HTTP/JSON REST requests to an API Gateway or gRPC-Gateway.

The gateway translates REST calls into gRPC calls and forwards them to internal services: User Service (port 50051), Order Service (port 50052), and Product Service (port 50053).

These internal services communicate with each other directly via gRPC over HTTP/2 with Protobuf payloads.

Services also publish events to a Message Queue (Kafka or RabbitMQ) for asynchronous processing.

This is the most common production pattern: REST at the edge, gRPC internally.

---

## Files in This Stage

None. This is a reading stage.

```
stage-02-grpc-vs-rest/
└── README.md   (you are here)
```

---

## How to Run This Stage

There is nothing to execute. Instead, read this README fully, reflect on your own experience with REST and whether you have encountered any of the problems listed, and think about a project you know and whether gRPC or REST would be better for it.

---

## What to Observe / Reflect On

After reading, you should be able to answer in your own words:

- What is the fundamental difference in how REST and gRPC model API calls?
- Why is Protobuf faster than JSON at parsing?
- Why can browsers not call gRPC directly?
- What does gRPC gain from HTTP/2 that REST (with HTTP/1.1) does not have?
- In what scenario would you use both REST and gRPC in the same system?
- What are two cases where REST is clearly the better choice?

---

## Exercise

There is no code exercise for this stage. Instead, think of an API you use or have built (for example, a user management API or a product catalog API). Write it out in REST style (GET, POST, etc.), then write the same thing in gRPC style (service plus methods). Finally, decide for that specific API whether you would prefer REST or gRPC, and write one to two sentences explaining why.

---

## What is Next

**Stage 03 — Client-Server Communication in gRPC** (still no code)

- Deep dive into the request lifecycle
- The stub, the channel, HTTP/2 frames
- What actually happens byte-by-byte

After you complete this stage, say **"ready for stage 03"** and we will continue.
