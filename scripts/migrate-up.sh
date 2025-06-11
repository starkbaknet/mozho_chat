#!/bin/bash

# Load environment variables
set -a
source .env
set +a

# Run migrate up
migrate -path ./migrations -database "$POSTGRES_URL" up
