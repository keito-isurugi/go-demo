# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a Go demo project containing multiple demonstration packages and examples for various Go programming concepts and libraries.

## Key Dependencies

- **GORM**: ORM library for database operations (PostgreSQL driver)
- **gqlgen**: GraphQL server library
- **Cobra**: CLI application framework
- **testify**: Testing assertions library

## Common Development Commands

### Database Operations
```bash
# Apply database schema (creates schema.sql from DDL/*.up.sql files)
make exec-schema

# Insert dummy data
make exec-dummy

# Refresh database (schema + dummy data)
make refresh-schema
```

### Testing
```bash
# Run all tests (note: some demo packages have their own go.mod files)
go test ./...

# Run specific test with verbose output
go test -v ./path/to/package

# Run tests with coverage
go test -cover ./...
```

### Running the Application
```bash
# Start the main HTTP server (runs on :8080)
go run main.go

# Build the application
go build -o go-demo main.go
```

## Project Architecture

### Main Application Structure
- **main.go**: HTTP server with endpoints for various demos
- **handler/**: HTTP request handlers
  - `algorithm_handler.go`: Algorithm demo endpoints
  - `time_handler.go`: Time-related demo endpoints
- **books/**: Book examples and design pattern implementations

### Demo Packages
Each demo package in `demo/` is self-contained with its own `go.mod`:
- **algorithm/**: Sorting and searching algorithms
- **crypto_demo/**: Cryptography examples (Caesar, Vigenere, DES, RSA, etc.)
- **designpattern/**: Design pattern implementations
- **oauth/**: OAuth implementation with tests
- **performance/**: Performance testing and benchmarks
- **zap/**: Structured logging examples

### Database
- Uses PostgreSQL with Docker container named `go-demo-db`
- DDL scripts in `DDL/` directory
- GORM models defined in `main.go`: Todo, User, Category, TodoCategory

### Testing Approach
- Test files follow Go convention: `*_test.go`
- Some packages use testify for assertions
- Mock generation for OAuth package in `demo/oauth/mock/`

## Important Notes

1. Many demo packages have their own `go.mod` files, making them independent modules
2. The main server exposes demo endpoints under `/demo/*` paths
3. Database connection uses hardcoded credentials in `main.go:35`
4. When adding new demos, consider whether they should be part of the main module or independent