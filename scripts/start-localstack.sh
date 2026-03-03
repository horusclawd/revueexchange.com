#!/bin/bash
set -e

echo "Starting LocalStack and local services..."

docker-compose up -d

echo "Waiting for LocalStack to be ready..."
max_attempts=30
attempt=0
until curl -s http://localhost:4566/_localstack/health | grep -q '"dynamodb": "available"'; do
  attempt=$((attempt + 1))
  if [ $attempt -ge $max_attempts ]; then
    echo "LocalStack did not become ready in time"
    docker-compose logs localstack
    exit 1
  fi
  echo "Waiting for LocalStack... ($attempt/$max_attempts)"
  sleep 5
done

echo ""
echo "LocalStack is ready!"
echo ""
echo "Services available at:"
echo "  - LocalStack: http://localhost:4566"
echo "  - PostgreSQL: localhost:5432 (user: revueadmin, pass: revueexchange, db: revueexchange)"
echo "  - Redis: localhost:6379"
echo ""
echo "Next steps:"
echo "  1. Deploy Terraform to LocalStack: make local-infra-deploy"
echo "  2. Start Go API: make local-api"
echo "  3. Start React UI: make local-ui"
