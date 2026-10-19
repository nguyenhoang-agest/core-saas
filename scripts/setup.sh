#!/usr/bin/env bash
set -euo pipefail

# ==============================================================================
# Core SaaS Monorepo — One-Click Development Setup
# ==============================================================================

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

GREEN="\033[32m"
CYAN="\033[36m"
YELLOW="\033[33m"
RED="\033[31m"
RESET="\033[0m"

echo -e "${CYAN}================================================================${RESET}"
echo -e "${CYAN}        Core SaaS — Initial Development Environment Setup       ${RESET}"
echo -e "${CYAN}================================================================${RESET}"
echo -e ""

cd "${ROOT_DIR}"

# 1. Environment Files
echo -e "${CYAN}==> [1/3] Generating environment configuration files...${RESET}"
if [ ! -f core-backend/.env ]; then
    cp core-backend/.env.example core-backend/.env
    echo -e "    ${GREEN}✓ Created core-backend/.env${RESET}"
else
    echo -e "    ${YELLOW}• core-backend/.env exists${RESET}"
fi

if [ ! -f core-frontend/.env.local ]; then
    cp core-frontend/.env.example core-frontend/.env.local
    echo -e "    ${GREEN}✓ Created core-frontend/.env.local${RESET}"
else
    echo -e "    ${YELLOW}• core-frontend/.env.local exists${RESET}"
fi

# 2. Check Docker and Start Database
echo -e "${CYAN}==> [2/3] Checking Docker & starting PostgreSQL container...${RESET}"
if ! docker info >/dev/null 2>&1; then
    echo -e "    ${RED}✗ Docker daemon is not running. Please start Docker and run 'make setup' again.${RESET}"
    exit 1
fi

docker compose -f core-backend/deploy/docker-compose.yml up -d

echo -e "    Waiting for PostgreSQL container to be healthy..."
RETRIES=15
until docker exec saas-starter-postgres pg_isready -U postgres -d saas_dev >/dev/null 2>&1 || [ $RETRIES -eq 0 ]; do
    sleep 1
    RETRIES=$((RETRIES - 1))
done

if [ $RETRIES -eq 0 ]; then
    echo -e "    ${RED}✗ PostgreSQL container failed to start in time.${RESET}"
    exit 1
fi
echo -e "    ${GREEN}✓ PostgreSQL container is healthy on port 5432 (DB: saas_dev).${RESET}"

# 3. Summary & Next Steps
echo -e ""
echo -e "${CYAN}==> [3/3] Ready for local development!${RESET}"
echo -e "${GREEN}================================================================${RESET}"
echo -e "${GREEN} Setup Complete! To start your development servers:${RESET}"
echo -e ""
echo -e "  Terminal 1 (Backend Monolith on :8080):"
echo -e "    ${CYAN}make dev-backend${RESET}   (or: cd core-backend && make run)"
echo -e ""
echo -e "  Terminal 2 (Frontend Admin UI on :3000):"
echo -e "    ${CYAN}make dev-frontend${RESET}  (or: cd core-frontend && npm run dev)"
echo -e ""
echo -e "  To run test suites across both backend and frontend:"
echo -e "    ${CYAN}make test${RESET}"
echo -e "${GREEN}================================================================${RESET}"
