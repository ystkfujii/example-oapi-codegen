# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is an OpenAPI code generation example using oapi-codegen for Go. The project demonstrates generating Go server and client code from an OpenAPI specification for a simple User API.

## Architecture

- **OpenAPI Spec**: `openapi/openapi.yaml` defines a User API with CRUD operations
- **Code Generation**: Uses oapi-codegen v2 to generate Go code from the spec
- **Generated Code**: `openapi/gen.go` contains generated models, client, and server interfaces
- **Server Implementation**: `server.go` implements the ServerInterface with an in-memory user store
- **Entry Point**: `main.go` sets up Echo server and registers handlers

The server uses:
- Echo framework for HTTP handling  
- Thread-safe in-memory storage with sync.RWMutex
- Auto-generated type-safe handlers from OpenAPI spec

## Development Commands

### Code Generation
```bash
go generate ./...
```
Regenerates Go code from the OpenAPI specification using oapi-codegen.config.yaml.

### Build and Run
```bash
go mod tidy    # Install dependencies
go run .       # Run the server on :8080
```

### Testing the API
```bash
# Get all users
curl http://localhost:8080/users

# Create a user  
curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{"name":"Alice","age":30}'

# Get user by ID
curl http://localhost:8080/users/1

# Delete user
curl -X DELETE http://localhost:8080/users/1
```

## Code Generation Configuration

The `oapi-codegen.config.yaml` configures code generation:
- Generates Echo server handlers and models
- Outputs to `openapi/gen.go` 
- Embeds the OpenAPI spec in generated code
- Client generation disabled (server-only)

When modifying the OpenAPI spec, run `go generate` to update generated code.