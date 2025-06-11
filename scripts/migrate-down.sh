#!/bin/bash

# Load environment variables
set -a
source .env
set +a

# Run migrate down
migrate -path ./migrations -database "$POSTGRES_URL" down
