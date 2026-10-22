#!/usr/bin/env bash
set -euo pipefail

# ==============================================================================
# Core SaaS Backend — Setup Environment & Database
# ==============================================================================

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(cd "${SCRIPT_DIR}/.." && pwd)"

GREEN="\033[32m"
CYAN="\033[36m"
YELLOW="\033[33m"
RED="\033[31m"
RESET="\033[0m"

echo -e "${CYAN}==> [1/3] Checking environment configuration...${RESET}"
cd "${ROOT_DIR}"

if [ ! -f .env ]; then
    cp .env.example .env
    echo -e "    ${GREEN}✓ Created core-backend/.env from .env.example${RESET}"
else
    echo -e "    ${YELLOW}• core-backend/.env already exists (preserving current settings)${RESET}"
fi

echo -e "${CYAN}==> [2/3] Checking Docker daemon status...${RESET}"
if ! docker info >/dev/null 2>&1; then
    echo -e "    ${RED}✗ Docker is not running. Please start Docker and run this script again.${RESET}"
    exit 1
fi
echo -e "    ${GREEN}✓ Docker daemon is active.${RESET}"

echo -e "${CYAN}==> [3/3] Starting local PostgreSQL 17 container...${RESET}"
docker compose -f deploy/docker-compose.yml up -d

echo -e "    Waiting for PostgreSQL to be healthy..."
RETRIES=15
until docker exec saas-starter-postgres pg_isready -U postgres -d saas_dev >/dev/null 2>&1 || [ $RETRIES -eq 0 ]; do
    sleep 1
    RETRIES=$((RETRIES - 1))
done

if [ $RETRIES -eq 0 ]; then
    echo -e "    ${RED}✗ Timed out waiting for PostgreSQL.${RESET}"
    exit 1
fi

echo -e "    ${GREEN}✓ PostgreSQL is healthy and accepting connections on localhost:5432 (database: saas_dev).${RESET}"
echo -e ""
echo -e "${GREEN}================================================================${RESET}"
echo -e "${GREEN}✓ Backend environment & database are ready!${RESET}"
echo -e "  To start the backend monolith:"
echo -e "    ${CYAN}cd core-backend && make run${RESET}  (or from root: ${CYAN}make dev-backend${RESET})"
echo -e "${GREEN}================================================================${RESET}"
