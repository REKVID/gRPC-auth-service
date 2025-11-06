
# Authentication Service

project demonstrating a microservices authentication system built with gRPC and HTTP Gateway. This project implements user registration and login functionality with JWT token-based authentication.

## Overview

This is an authentication service that showcases Go backend architecture patterns. The service provides both gRPC and REST API interfaces through a unified protocol buffer definition, allowing clients to interact via either protocol.

## Technologies and Stack

### Core Technologies

- **Go 1.24+** - Primary programming language
- **gRPC** - High-performance RPC framework for service-to-service communication
- **Protocol Buffers** - Language-neutral data serialization format for API contracts
- **gRPC Gateway** - HTTP/REST to gRPC translation layer

### Libraries and Frameworks

- **google.golang.org/grpc** - gRPC implementation for Go
- **github.com/grpc-ecosystem/grpc-gateway/v2** - Generates HTTP reverse proxy from proto definitions
- **google.golang.org/protobuf** - Protocol Buffers support for Go
- **gorm.io/gorm** - ORM for database operations
- **gorm.io/driver/postgres** - PostgreSQL driver for GORM
- **github.com/golang-jwt/jwt/v5** - JWT token generation and validation
- **golang.org/x/crypto/bcrypt** - Password hashing using bcrypt algorithm

### Database

- **PostgreSQL** - Relational database for user data storage

## Architecture

### Service Architecture

The project follows a layered architecture pattern with clear separation of concerns:

1. **API Layer** (`api/proto/`) - Protocol buffer definitions and generated code
2. **Handler Layer** (`internal/auth/handler.go`) - gRPC request handlers
3. **Service Layer** (`internal/auth/service.go`) - Business logic implementation
4. **Data Layer** (`internal/store/`) - Database access and data persistence

### Application Structure

```
auth2/
├── api/
│   └── proto/
│       ├── auth.proto              # Protocol buffer service definition
│       ├── auth.pb.go              # Generated message types
│       ├── auth_grpc.pb.go         # Generated gRPC service code
│       ├── auth.pb.gw.go           # Generated HTTP gateway code
│       └── google/api/             # Google API proto definitions
├── cmd/
│   ├── server/                     # gRPC server application
│   │   └── main.go
│   └── gateway/                    # HTTP gateway application
│       └── main.go
├── internal/
│   ├── auth/                       # Authentication domain logic
│   │   ├── handler.go              # gRPC handler implementation
│   │   └── service.go              # Business logic (register, login)
│   └── store/                      # Data persistence layer
│       └── store.go                # Database operations
├── go.mod                          # Go module dependencies
└── go.sum                          # Dependency checksums
```

### Design Principles

1. **Separation of Concerns** - Each layer has a single responsibility
2. **Dependency Injection** - Services and handlers receive dependencies through constructors
3. **Interface-Based Design** - gRPC handlers implement generated service interfaces
4. **Protocol-First Development** - API contracts defined in proto files before implementation
5. **Code Generation** - Client and server code generated from proto definitions

### Communication Flow

1. **HTTP Request Flow**: Client → HTTP Gateway → gRPC Server → Service → Store → Database
2. **gRPC Request Flow**: Client → gRPC Server → Service → Store → Database

The HTTP Gateway translates REST requests to gRPC calls, allowing the same backend to serve both protocols.

## Features

### Authentication Features

- **User Registration** - Create new user accounts with email and password
- **User Login** - Authenticate users and receive JWT tokens
- **Password Security** - Passwords hashed using bcrypt algorithm
- **JWT Tokens** - Stateless authentication tokens for session management

### API Capabilities

- **Dual Protocol Support** - Access via gRPC or HTTP/REST
- **Type-Safe Contracts** - Protocol buffers ensure consistent data structures
- **Automatic Code Generation** - Client and server stubs generated from proto files

## Prerequisites

- Go 1.24 or higher
- PostgreSQL 12 or higher
- Protocol Buffer compiler (protoc)
- Go plugins for protoc:
  - protoc-gen-go
  - protoc-gen-go-grpc
  - protoc-gen-grpc-gateway

## Installation

### Install Protocol Buffer Compiler

**macOS:**
```bash
brew install protobuf
```

**Linux:**
```bash
sudo apt-get install protobuf-compiler
```

**Or download from:** https://github.com/protocolbuffers/protobuf/releases

### Install Go Protobuf Plugins

```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
```

Add Go bin directory to PATH:
```bash
export PATH="$PATH:$(go env GOPATH)/bin"
```

### Setup PostgreSQL Database

Create database and user:
```sql
CREATE DATABASE authdb;
CREATE USER auth WITH PASSWORD 'auth';
GRANT ALL PRIVILEGES ON DATABASE authdb TO auth;
```

## Code Generation

Generate Go code from protocol buffer definitions:

```bash
protoc -I api/proto \
  --go_out=. --go_opt=paths=source_relative \
  --go-grpc_out=. --go-grpc_opt=paths=source_relative \
  --grpc-gateway_out=. --grpc-gateway_opt=paths=source_relative \
  api/proto/auth.proto
```

This generates three files in `api/proto/`:
- `auth.pb.go` - Message type definitions
- `auth_grpc.pb.go` - gRPC service interface
- `auth.pb.gw.go` - HTTP gateway handlers

## Running the Application

### Start gRPC Server

The gRPC server listens on port 50051:

```bash
go run cmd/server/main.go
```

Server output:
```
gRPC server listening on :50051
```

### Start HTTP Gateway

The HTTP gateway listens on port 8080 and forwards requests to the gRPC server:

```bash
go run cmd/gateway/main.go
```

Gateway output:
```
HTTP Gateway running on :8080
```

**Note:** The gateway requires the gRPC server to be running on `localhost:50051`.

## API Endpoints

### HTTP REST API

Base URL: `http://localhost:8080`

#### Register User

```http
POST /v1/auth/register
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "securepassword123"
}
```

Response:
```json
{
  "token": "user_1"
}
```

#### Login User

```http
POST /v1/auth/login
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "securepassword123"
}
```

Response:
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

### gRPC API

Connect to: `localhost:50051`

#### Register RPC

```protobuf
rpc Register(RegisterRequest) returns (RegisterResponse)
```

Request:
```protobuf
message RegisterRequest {
  string email = 1;
  string password = 2;
}
```

#### Login RPC

```protobuf
rpc Login(LoginRequest) returns (LoginResponse)
```

Request:
```protobuf
message LoginRequest {
  string email = 1;
  string password = 2;
}
```

## Configuration

### Database Connection

Edit `cmd/server/main.go` to change database connection string:

```go
store.NewStore("postgres://user:password@host:port/database?sslmode=disable")
```

### JWT Secret Key

Edit `cmd/server/main.go` to change JWT signing key:

```go
auth.NewService(store, "your-secret-key-here")
```

### Server Ports

- gRPC Server: Port 50051 (edit in `cmd/server/main.go`)
- HTTP Gateway: Port 8080 (edit in `cmd/gateway/main.go`)

## Testing the API

### Using curl

Register a user:
```bash
curl -X POST http://localhost:8080/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"test123"}'
```

Login:
```bash
curl -X POST http://localhost:8080/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"test123"}'
```

### Using gRPC Client

You can use tools like `grpcurl` or generate a client from the proto file to test gRPC endpoints directly.

## Project Structure Explanation

### `api/proto/`

Contains protocol buffer definitions and generated code. The proto file defines the service contract, and code generation creates type-safe Go code.

### `cmd/server/`

gRPC server application entry point. Initializes database connection, creates service layer, and registers gRPC handlers.

### `cmd/gateway/`

HTTP gateway application. Translates HTTP requests to gRPC calls and forwards them to the gRPC server.

### `internal/auth/`

Authentication domain logic:
- `service.go` - Business logic for registration and login
- `handler.go` - gRPC request handlers that call service methods

### `internal/store/`

Data persistence layer:
- Database connection management
- User model definition
- CRUD operations for user data

## Security Considerations

This is a learning project. For production use, consider:

- Use environment variables for sensitive configuration
- Implement proper JWT token expiration and refresh
- Add rate limiting to prevent brute force attacks
- Use TLS/SSL for all connections
- Implement proper error handling and logging
- Add input validation and sanitization
- Use connection pooling for database
- Implement proper password strength requirements

## Development Workflow

1. Modify `api/proto/auth.proto` to change API contract
2. Regenerate code using `protoc` command
3. Update service and handler implementations
4. Test with HTTP or gRPC clients
5. Restart servers to apply changes

## Dependencies

All dependencies are managed through Go modules. Install with:

```bash
go mod download
```
