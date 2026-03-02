#!/bin/bash
set -e

echo "Stopping LocalStack and local services..."

docker-compose down

echo ""
echo "All services stopped."
echo "Data preserved in Docker volumes."
