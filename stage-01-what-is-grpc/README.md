# Stage 01 — What is gRPC and Why It Exists

> **Phase**: Foundations (Read Only)
> **Type**: Concept-only stage — no code, no Go, no files to run
> **Goal**: Build the mental model before writing a single line of code

---

## 🎯 What I'm Learning

In this stage, you will understand:

1. What "Remote Procedure Call" (RPC) actually means
2. What gRPC is and where it came from
3. The three pillars that make gRPC different: **Protocol Buffers**, **HTTP/2**, and **contract-first design**
4. Why Google built it and what problem it solves

---

## 🤔 Why This Exists

Before gRPC, Google had thousands of services talking to each other across data centers. They used a system internally called **Stubby** (built in 2001). It worked, but:

- It only worked inside Google (not open source)
- It didn't support many languages
- It wasn't built on modern web standards

In **2015**, Google open-sourced a complete redesign as **gRPC** and donated it to the **Cloud Native Computing Foundation (CNCF)**. Today, gRPC is the de-facto standard for service-to-service communication in microservices.

---

## 🧱 The Three Pillars of gRPC

### 1. Protocol Buffers (Protobuf) — The Data Format

When two programs talk to each other, they need to agree on how data is encoded.

| Format | Human Readable | Size | Speed | Used By |
|---|---|---|---|---|
| JSON | ✅ Yes | Large | Slow | REST APIs |
| XML | ✅ Yes | Very large | Very slow | Old SOAP APIs |
| **Protobuf** | ❌ No (binary) | **Small** | **Very fast** | **gRPC** |

Protobuf is a **binary serialization format**. Instead of sending:
```json
{"name":"Alice","age":30}
```
Protobuf sends a compact binary representation that is much smaller and faster to parse.

But Protobuf alone isn't useful — we need a way to **describe** what data looks like. That's where `.proto` files come in.

### 2. HTTP/2 — The Transport Protocol

gRPC runs on **HTTP/2**, not HTTP/1.1. This is a major reason it's so fast:

| Feature | HTTP/1.1 | HTTP/2 |
|---|---|---|
| Multiple requests per connection | ❌ No (one at a time) | ✅ Yes (multiplexed) |
| Header compression | ❌ No | ✅ Yes (HPACK) |
| Binary framing | ❌ No (text) | ✅ Yes (binary) |
| Server push | ❌ No | ✅ Yes |

This means a single TCP connection can carry many gRPC calls simultaneously, with tiny overhead.

### 3. Contract-First Design (`.proto` files)

In REST, you often write code first and document later. In gRPC, you write the **contract first**:

```
┌────────────────────┐
│  greeting.proto    │  ← Contract: "there is a Greet method that takes a name
└────────┬───────────┘              and returns a greeting"
         │
         │  protoc compiles it
         │
         ├─────────────────────────┐
         ▼                         ▼
┌─────────────────┐      ┌──────────────────┐
│  Go client      │      │  Python client   │  ← Both generated from the
│  (c.Greet())    │      │  (c.greet())    │     same .proto file
└─────────────────┘      └──────────────────┘
```

---

## 🔄 How gRPC Works (The Big Picture)

```
┌──────────────┐                              ┌──────────────┐
│              │   1. Call c.Greet(name)      │              │
│              ├─────────────────────────────►│              │
│   Go Client  │                              │  Go Server   │
│              │   2. Return response         │              │
│              │◄─────────────────────────────┤              │
└──────────────┘                              └──────────────┘

        │  What's actually happening under the hood:
        │  ┌─────────────────────────────────────────────────┐
        │  │ 1. Marshal GreetRequest to binary Protobuf    │
        │  │ 2. Send over HTTP/2 (binary frames)          │
        │  │ 3. Server unmarshals binary                   │
        │  │ 4. Calls your Greet() handler                │
        │  │ 5. Marshals response to binary                │
        │  │ 6. Sends back over HTTP/2                    │
        │  │ 7. Client unmarshals response                 │
        │  └─────────────────────────────────────────────────┘
```

To the developer, it looks like a normal function call. But underneath, there's serialization, network transport, and deserialization happening automatically.

---

## 📦 What's in a gRPC "Call"?

Every gRPC call has 4 parts:

1. **Service** — the API endpoint (e.g., `GreetService`)
2. **Method** — the function to call (e.g., `Greet`)
3. **Request message** — input data (defined in `.proto`)
4. **Response message** — output data (defined in `.proto`)

You'll write all 4 in your `.proto` file.

---
---

## 🆚 Quick Glimpse: gRPC vs REST

We'll go deep in Stage 02, but here's a taste:

| | REST | gRPC |
|---|---|---|
| Format | JSON (text) | Protobuf (binary) |
| Transport | HTTP/1.1 or HTTP/2 | HTTP/2 only |
| Contract | Optional (OpenAPI) | Required (`.proto`) |
| Browser support | Native | Needs gateway |
| Streaming | Not built-in | ✅ Built-in |
| Speed | Slower | Faster |

---

## 📂 Files in This Stage

**None.** This is a reading stage.

```
stage-01-what-is-grpc/
└── README.md   ← (you are here)
```

---

## ▶️ How to "Run" This Stage

There's nothing to execute. Instead:

1. **Read** this README fully
2. **Read** the official gRPC introduction: https://grpc.io/docs/what-is-grpc/introduction/
3. **Skim** the Protobuf overview: https://protobuf.dev/overview/
4. **Watch** the request flow diagram above — make sure you understand each step

---

## 👀 What to Observe / Reflect On

After reading, you should be able to answer these in your own words:

- [ ] What does "RPC" mean? How is it different from a normal function call?
- [ ] Why did Google build gRPC? What was wrong with their previous solution?
- [ ] What are the three pillars of gRPC and what does each one do?
- [ ] Why is binary (Protobuf) faster than text (JSON)?
- [ ] Why does gRPC use HTTP/2 instead of HTTP/1.1?
- [ ] What is a "contract-first" API design?

If any of these are unclear, ask questions before moving on. The next stages build directly on this mental model.

---

## ✏️ Exercise

There is no code exercise for this stage. Instead, do this:

1. **Write down** in your own words (a notebook, a text file, anywhere):
   - What gRPC is
   - Why it exists
   - The three pillars

2. **Draw** the request flow diagram from memory. You should be able to draw:
   - Client calling a method
   - The stub marshaling the request
   - HTTP/2 sending it to the server
   - The server unmarshaling and calling your handler
   - The response flowing back

3. **Answer** the self-check questions above without looking at this README.

---

## ⏭️ What's Next

**Stage 02 — gRPC vs REST** (still no code)
- Deep comparison of the two
- When to use which
- Tradeoffs in real systems

After you complete this stage, say **"ready for stage 02"** and we'll continue.
