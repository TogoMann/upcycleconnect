#!/bin/bash

# UpcycleConnect Test Runner
# This script runs both backend and frontend tests.

set -e

# ANSI Color Codes
GREEN='\033[0;32m'
BLUE='\033[0;34m'
RED='\033[0;31m'
NC='\033[0m'

echo -e "${BLUE} ---------- Starting tests ----------${NC}"

if [ -f .env ]; then
    set -a
    source .env
    set +a
fi

# 1. Run Backend Tests (Go)
echo -e "${GREEN}[1/2] Testing backend...${NC}"
cd backend
if go test ./tests/... -v; then
    echo -e "${GREEN}✓ Backend tests passed.${NC}\n"
else
    echo -e "${RED}✗ Backend tests failed.${NC}\n"
    exit 1
fi
cd ..

# 2. Run Frontend Tests (Vitest via Bun in Docker)
echo -e "${GREEN}[2/2] Testing frontend...${NC}"
# Check if the frontend container is running
if ! docker ps --format '{{.Names}}' | grep -q "^uc-front-dev$"; then
    echo -e "${RED}✗ Error: Docker container 'uc-front-dev' is not running.${NC}"
    echo -e "Please start your development environment before running frontend tests."
    exit 1
fi

if docker exec uc-front-dev bun run test:unit -- tests/ --run; then
    echo -e "\n${GREEN}✓ Frontend tests passed.${NC}\n"
else
    echo -e "\n${RED}✗ Frontend tests failed.${NC}\n"
    exit 1
fi

echo -e "${GREEN}---------- All tests ok ----------${NC}"
