# Stage 04 — Writing Your First .proto File

> **Phase**: Protocol Buffers (Hands-On)
> **Type**: Code stage — you will write a `.proto` file
> **Goal**: Write your first `.proto` file from scratch and understand every part of it
> **Prerequisites**: Stage 01, Stage 02, Stage 03

---

## What I am Learning

1. What a `.proto` file is and why it is the "contract" between client and server
2. What Protocol Buffers (Protobuf) are and why they are different from JSON
3. The anatomy of a `.proto` file — every keyword explained
4. What field numbers are and why they matter
5. How `protoc` turns a `.proto` file into Go code

---

## What is New vs Stage 03

Stage 03 was purely conceptual — we traced the request flow. This stage is the first **hands-on stage**. You will write real files.

---

## Why .proto Files Exist

A `.proto` file is the **contract** between your client and server. Both sides agree on it before writing any code.

```
greet.proto (the contract)
       │
       │  protoc compiles it
       ▼
┌──────────────┐     ┌──────────────┐
│  Go server   │     │  Go client  │
│  implements  │     │  calls the  │
│  the methods │     │  stub       │
└──────────────┘     └──────────────┘
```

Both the server and client are **generated from the same file**. This guarantees they always agree on:
- What methods exist
- What fields are in each message
- What types each field is

If you change the `.proto` file and regenerate, the compiler tells you immediately if something is wrong — not at runtime.

---

## The Mechanism: .proto → protoc → Go Code

```
┌─────────────────┐
│  greet.proto    │  ← You write this (the contract)
│  (plain text)   │
└────────┬────────┘
         │ protoc (Protocol Buffer compiler)
         │ + protoc-gen-go plugin
         │ + protoc-gen-go-grpc plugin
         ▼
┌─────────────────────────────────┐
│  greet.pb.go     │  greet_grpc.pb.go
│  (message types) │  (client stub + server interface)
│  DO NOT EDIT     │  DO NOT EDIT
└─────────────────────────────────┘
```

`greet.pb.go` contains:
- Go structs for every message (e.g., `GreetRequest`, `GreetResponse`)
- Marshal and unmarshal methods

`greet_grpc.pb.go` contains:
- `GreetServiceClient` — the interface the client uses
- `NewGreetServiceClient` — constructor for the client
- `GreetServiceServer` — the interface your server must implement
- `RegisterGreetServiceServer` — to register your server with gRPC

---

## The Anatomy of a .proto File

Here is a complete, minimal `.proto` file. Every part is explained below.

```protobuf
// File: proto/greet.proto

// 1. Syntax declaration — always proto3
syntax = "proto3";

// 2. Package — groups related files, like a namespace
package greet;

// 3. Go package option — tells protoc where to put generated Go files
option go_package = "stage-05-generating-go-code/gen/greet;greetpb";

// 4. Message definitions — the data structures

message GreetRequest {
    // field: type = number
    string name = 1;
}

message GreetResponse {
    string result = 1;
}

// 5. Service definition — the RPC methods

service GreetService {
    // rpc MethodName(RequestType) returns (ResponseType);
    rpc Greet(GreetRequest) returns (GreetResponse);
}
```

---

## Understanding Every Part

### `syntax = "proto3";`

Always the first line. There are two versions: `proto2` and `proto3`. Always use `proto3` for new projects. It is simpler and supports more languages.

### `package greet;`

A namespace for this file. Similar to a Go package. Prevents name conflicts between different `.proto` files.

### `option go_package = "...";`

This tells protoc where to put the generated Go files. It has two parts separated by a semicolon:
- The first part: the directory path (e.g., `stage-05-generating-go-code/gen/greet`)
- The second part: the Go package name (e.g., `greetpb`)

You will import the package using the second part: `import "stage-05-generating-go-code/gen/greet;greetpb"`

### Message Definitions

```protobuf
message GreetRequest {
    string name = 1;
}
```

A **message** is a struct. It has **fields**. Each field has:
- A **type** (`string`, `int32`, `bool`, etc.)
- A **name** (`name`)
- A **field number** (`1`)

### Field Numbers — The Most Important Concept

Field numbers are what makes Protobuf work. Here is why they matter:

In JSON, you send the field **name**:
```json
{"name": "Alice", "age": 30}
```

In Protobuf, you send the field **number**:
```
[name=1, value="Alice"], [name=2, value=30]
```

Or more precisely, the wire format uses the number directly (not even "name=1" — just `1` is the tag). This is why Protobuf is so compact.

**Rules for field numbers:**
- Must be **positive integers**
- Must be **unique** within a message
- Numbers **1-15** use 1 byte on the wire (use for frequently-used fields)
- Numbers **16-2047** use 2+ bytes
- **Never change** a field number once the file is in use — old clients/servers would break
- You **can add** new fields with new numbers — old code ignores them

---

### Field Types

Scalar types in Protobuf:

| Proto Type | Go Type | Description |
|---|---|---|
| `string` | `string` | UTF-8 text |
| `int32` | `int32` | Signed 32-bit integer |
| `int64` | `int64` | Signed 64-bit integer |
| `uint32` | `uint32` | Unsigned 32-bit integer |
| `bool` | `bool` | True or false |
| `double` | `float64` | 64-bit float |
| `float` | `float32` | 32-bit float |
| `bytes` | `[]byte` | Raw bytes |

### Service and RPC Method

```protobuf
service GreetService {
    rpc Greet(GreetRequest) returns (GreetResponse);
}
```

A **service** groups related RPC methods. Each method has:
- A **name** (`Greet`)
- An **input type** (`GreetRequest`)
- An **output type** (`GreetResponse`)

Both the input and output must be **messages** (defined elsewhere in the file), not raw types.

---

## The .proto File You Will Write

A reference version of this file has been created at `proto/greet.proto` so you can see what it looks like. For this stage, **type it out yourself** (or copy and modify) to get hands-on practice:

```protobuf
// proto/greet.proto
syntax = "proto3";

package greet;

option go_package = "stage-05-generating-go-code/gen/greet;greetpb";

message GreetRequest {
    string name = 1;
}

message GreetResponse {
    string result = 1;
}

service GreetService {
    rpc Greet(GreetRequest) returns (GreetResponse);
}
```

This is the **smallest possible gRPC contract**. It defines:
- A request message with one field
- A response message with one field
- One service with one unary RPC method

---

## Files in This Stage

```
stage-04-first-proto-file/
├── README.md        ← (you are here)
├── proto/
│   └── greet.proto  ← the contract you will write
└── gen/             ← generated Go code will go here in Stage 05
```

---

## How to Run This Stage

### Step 1: Write the .proto file

Create `proto/greet.proto` with the content shown above.

### Step 2: Verify it is valid (optional)

If you have `protoc` installed, you can check the file syntax. But we will not generate Go code in this stage — we do that in Stage 05.

### Step 3: Read the .proto file again

After writing it, read the file line by line and confirm you can explain:
- What each `syntax`, `package`, `option`, and `message` keyword does
- Why `name = 1` uses the number `1` and not `0` or `2`
- What would happen if two fields had the same number
- What would happen if you removed `option go_package`

---

## What to Observe / Reflect On

After completing this stage, you should be able to explain:

- [ ] What is the purpose of the `go_package` option?
- [ ] Why are field numbers important in Protobuf?
- [ ] What happens if you assign the same field number to two fields in the same message?
- [ ] What is the difference between a `message` and a `service`?
- [ ] Why do RPC method parameters and return types need to be message types, not raw types?
- [ ] What does `protoc` produce when it compiles a `.proto` file?

---

## Exercise

Modify the `.proto` file to add a new RPC method called `GreetWithAge` that takes a `GreetRequest` (which already has `name`) and returns a `GreetResponse` that says `"Hello, <name>! You are <age> years old."`.

You will need to:
1. Add an `int32 age = 2;` field to `GreetRequest`
2. Keep the existing `name` field (do not change its number — `name = 1` stays)
3. Add a new RPC method `GreetWithAge` that also takes `GreetRequest` and returns `GreetResponse`

Do NOT generate Go code yet. Just write the modified `.proto` file.

---

## What is Next

**Stage 05 — Installing Tools and Generating Go Code**

- Install `protoc` (the Protocol Buffer compiler)
- Install `protoc-gen-go` and `protoc-gen-go-grpc` (the Go plugins)
- Run `protoc` to generate Go code from your `.proto` file
- Read the generated `.pb.go` files line by line
- Verify the structure matches what you wrote in the proto

After you complete this stage, say **"ready for stage 05"** and we will continue.
