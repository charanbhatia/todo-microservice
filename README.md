# Auth + Todo Microservice

A Go kit-based microservice that provides authentication and todo management functionality.

## Features

- **Authentication Service**: User signup, login, and token validation
- **Todo Service**: CRUD operations for user-owned todos
- **Service Middleware**: Logging and Prometheus metrics
- **Endpoint Middleware**: Rate limiting on write operations
- **In-Memory Storage**: Simple map-based persistence (ready for DB integration)

## Architecture

Built following Go kit patterns:

- **Service Layer**: Pure business logic interfaces and implementations
- **Endpoints**: RPC-style request/response adapters
- **HTTP Transport**: JSON over HTTP handlers
- **Middlewares**: Logging, metrics, and rate limiting

## Getting Started

### Prerequisites

- Go 1.21 or higher

### Installation

```bash
go mod download
```

### Build

```bash
go build -o bin/server.exe .
```

### Run

```bash
./bin/server.exe
```

Server starts on `http://localhost:8080`

## API Endpoints

### Authentication

**Signup**
```bash
curl -X POST http://localhost:8080/signup \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"pass123"}'
```

**Login**
```bash
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"pass123"}'
```

**Validate Token**
```bash
curl -X POST http://localhost:8080/validate \
  -H "Content-Type: application/json" \
  -d '{"token":"YOUR_TOKEN"}'
```

### Todo Management

**Create Todo**
```bash
curl -X POST http://localhost:8080/todos \
  -H "Content-Type: application/json" \
  -d '{"user_id":"user_1","text":"Buy groceries"}'
```

**List Todos**
```bash
curl -X GET "http://localhost:8080/todos?user_id=user_1"
```

**Complete Todo**
```bash
curl -X POST http://localhost:8080/todos/todo_1/complete \
  -H "Content-Type: application/json" \
  -d '{"user_id":"user_1"}'
```

### Metrics

**Prometheus Metrics**
```bash
curl http://localhost:8080/metrics
```

## Middleware

### Service Middleware

- **Logging**: Logs method calls with parameters, latency, and errors
- **Instrumentation**: Prometheus metrics for request count and latency

### Endpoint Middleware

- **Rate Limiting**: Token bucket rate limiter (10 req/s, burst 20) on:
  - Signup endpoint
  - Create todo endpoint

## Project Structure

```
.
├── auth_todo/
│   ├── service.go          # Service interfaces and implementations
│   ├── service_test.go     # Unit tests
│   ├── endpoints.go        # Endpoint definitions
│   ├── transport_http.go   # HTTP handlers and decoders
│   ├── middleware.go       # Logging and metrics middleware
│   └── ratelimit.go        # Rate limiting middleware
├── main.go                 # Application entry point
├── go.mod
└── README.md
```

## Metrics Available

- `auth_todo_auth_service_request_count`: Auth service request counter
- `auth_todo_auth_service_request_latency_microseconds`: Auth service latency
- `auth_todo_todo_service_request_count`: Todo service request counter
- `auth_todo_todo_service_request_latency_microseconds`: Todo service latency

## Kubernetes Deployment

### Build Docker Image

```bash
docker build -t todo-microservice:latest .
```

### Deploy to Kubernetes

```bash
kubectl apply -f k8s/namespace.yaml
kubectl apply -f k8s/deployment.yaml
kubectl apply -f k8s/service.yaml
```

### Check Deployment

```bash
kubectl get pods -n todo-microservice
kubectl get svc -n todo-microservice
```

### Access Service

For local clusters (minikube):
```bash
minikube service todo-microservice -n todo-microservice
```

For cloud providers, get the external IP:
```bash
kubectl get svc todo-microservice -n todo-microservice
```

### Scale Deployment

```bash
kubectl scale deployment todo-microservice --replicas=3 -n todo-microservice
```

## Kubernetes Configuration

- **Deployment**: 2 replicas with liveness and readiness probes
- **Service**: LoadBalancer type exposing port 80
- **Resources**: 64Mi-128Mi memory, 100m-200m CPU per pod
- **Namespace**: Isolated namespace for the service

## Future Enhancements

- Replace in-memory storage with PostgreSQL or Redis
- Add JWT signing with proper library
- Implement distributed tracing
- Add service discovery and client-side load balancing
- Create CLI client using Go kit
