# Traefik Test with ForwardAuth

A demonstration of Traefik v3 with ForwardAuth middleware using Go services and Redis.

## Architecture

- **Traefik**: Reverse proxy with ForwardAuth middleware
- **Auth Service**: Handles API key validation and management (Go with standard library HTTP router)
- **API Service**: Protected and public endpoints (Go with standard library HTTP router)
- **Redis**: API key storage

## Quick Start

### Prerequisites

- Docker and Docker Compose
- curl (for testing)

### 1. Start the Services

```bash
docker compose up -d
```

This will start:
- Traefik on port 80 (main) and 8080 (dashboard)
- Redis on port 6379
- Auth service (internal)
- API service (internal)

### 2. Verify Services are Running

```bash
docker compose ps
```

### 3. Check Health

```bash
curl http://localhost/auth/health
```

Expected response: `Auth service healthy`

## API Key Management

### Create an API Key

```bash
curl -X POST http://localhost/auth/create-key \
  -d "user_id=john_doe&api_key=secret123"
```

Expected response: `API key created for user: john_doe`

### Test API Access

Without authentication (should fail):
```bash
curl http://localhost/api
```

Expected response: `Missing Authorization header`

With valid API key:
```bash
curl -H "Authorization: Bearer secret123" http://localhost/api/protected
```

The request will be forwarded to the API service with user headers.

Test public endpoints (no auth required):
```bash
curl http://localhost/api/health
curl http://localhost/api/public
```

## Endpoints

### Auth Service (`/auth`)

- `GET /auth/health` - Health check
- `POST /auth/create-key` - Create API key (requires `user_id` and `api_key` form data)
- `GET/POST /auth` - Internal ForwardAuth endpoint (used by Traefik)

### API Service (`/api`)

#### Protected Endpoints (require authentication):
- `GET /api/protected` - Protected endpoint with user info
- `GET /api/data` - Protected data endpoint
- Requires `Authorization: Bearer <api-key>` header

#### Public Endpoints (no authentication required):
- `GET /api/health` - Health check
- `GET /api/public` - Public endpoint

## Development

### Technology Stack

- **Go Services**: Built with Go 1.24+ using the standard library HTTP router (zero external dependencies for routing)
- **Redis**: For API key storage
- **Docker**: Containerized deployment
- **Traefik**: Advanced routing with separate routers for protected vs public endpoints

### Rebuild Services

After making code changes:

```bash
# Rebuild specific service
docker compose build auth-service
docker compose up -d auth-service

# Or rebuild all services
docker compose build
docker compose up -d
```

### View Logs

```bash
# All services
docker compose logs -f

# Specific service
docker compose logs -f traefik
docker compose logs -f auth-service
```

### Traefik Dashboard

Visit http://localhost:8080 to see the Traefik dashboard with routing information.

## Cleanup

```bash
docker compose down
```

To remove volumes as well:
```bash
docker compose down -v
```