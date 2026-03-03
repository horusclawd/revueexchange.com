# Local Development

This guide covers setting up local development for RevUExchange.

## Prerequisites

- Docker & Docker Compose
- Go 1.21+
- Node.js 18+

## Quick Start

### 1. Start Local Infrastructure

```bash
make local-start
```

This starts:
- LocalStack (AWS mock)
- PostgreSQL (database)
- Redis (caching)

### 2. Run Database Migration

The database needs to be initialized with the schema. Run:

```bash
docker exec -i revueexchange-postgres psql -U revueadmin -d revueexchange < api/migrations/001_initial_schema.sql
```

### 3. Create Environment File

Create `api/.env` with these variables:

```bash
# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=revueadmin
DB_PASSWORD=revueexchange
DB_NAME=revueexchange

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379

# AWS (LocalStack)
AWS_REGION=us-east-1
AWS_ENDPOINT=http://localhost:4566

# JWT
JWT_SECRET=your-secret-key-change-in-production
```

### 4. Start Development Servers

**Go API:**
```bash
make local-api
```

**React UI:**
```bash
make local-ui
```

## Services

| Service | URL | Credentials |
|--------|-----|-------------|
| LocalStack | http://localhost:4566 | - |
| PostgreSQL | localhost:5432 | revueadmin / revueexchange |
| Redis | localhost:6379 | - |
| Go API | http://localhost:8080 | - |
| React UI | http://localhost:5173 | - |

## Troubleshooting

### LocalStack not starting
```bash
docker compose logs localstack
```

### Database connection issues
```bash
docker compose restart postgres
```

### Clear all data
```bash
docker compose down -v
```
