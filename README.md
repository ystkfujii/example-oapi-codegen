# Example oapi-codegen

A sample Go HTTP server project auto-generated from OpenAPI specifications.

## Overview

This project demonstrates how to automatically generate server code from OpenAPI 3.0 specifications using oapi-codegen and implement a REST API server with the Echo framework.

## Usage

### Starting the Server

```bash
# Run with Docker Compose
make docker-compose-up
```

The server will start at `http://localhost:8080`.

### API Usage Examples

```bash
# Create a user
curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{"name":{"first": "bob", "last":"johnson"},"age":10}'

# Retrieve all users
curl http://localhost:8080/users

# Retrieve a specific user
curl http://localhost:8080/users/1

# Delete a user
curl -X DELETE http://localhost:8080/users/1
```
