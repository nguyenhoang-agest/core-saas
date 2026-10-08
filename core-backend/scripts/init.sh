#!/usr/bin/env bash
set -euo pipefail

# Usage: ./scripts/init.sh [PROJECT=<name>] [MODULE=<path>] [KEEP=<mod1,mod2>]
# Example: ./scripts/init.sh PROJECT=my-saas MODULE=github.com/my-org/my-saas KEEP=iam,tenant

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(cd "${SCRIPT_DIR}/.." && pwd)"

# Parse key=value arguments
PROJECT_NAME=""
NEW_MODULE=""
KEEP_MODULES=""

for arg in "$@"; do
    case "$arg" in
        PROJECT=*)
            PROJECT_NAME="${arg#*=}"
            ;;
        MODULE=*)
            NEW_MODULE="${arg#*=}"
            ;;
        KEEP=*)
            KEEP_MODULES="${arg#*=}"
            ;;
        *)
            echo "Unknown argument: $arg"
            echo "Usage: $0 [PROJECT=<name>] [MODULE=<path>] [KEEP=<mod1,mod2>]"
            exit 1
            ;;
    esac
done

cd "${ROOT_DIR}"

OLD_MODULE=$(grep -E '^module ' go.mod | awk '{print $2}')
if [ -z "${OLD_MODULE}" ]; then
    echo "Error: unable to determine current Go module path from go.mod"
    exit 1
fi

sed_i() {
    if [[ "$OSTYPE" == "darwin"* ]]; then
        sed -i '' "$@"
    else
        sed -i "$@"
    fi
}

# 1. Rename module path if specified
if [ -n "${NEW_MODULE}" ] && [ "${NEW_MODULE}" != "${OLD_MODULE}" ]; then
    echo "Renaming Go module: '${OLD_MODULE}' -> '${NEW_MODULE}'..."
    find . -type f \( -name "*.go" -o -name "go.mod" -o -name "*.md" -o -name "Makefile" -o -name "*.sh" \) \
        -not -path "*/.git/*" | while read -r f; do
        sed_i "s|${OLD_MODULE}|${NEW_MODULE}|g" "$f"
    done
fi

# 2. Rename project name if specified
if [ -n "${PROJECT_NAME}" ]; then
    echo "Updating project references to '${PROJECT_NAME}'..."
    find . -type f \( -name "*.md" -o -name "*.yaml" -o -name "*.yml" -o -name "Dockerfile" \) \
        -not -path "*/.git/*" | while read -r f; do
        sed_i "s|core-backend|${PROJECT_NAME}|g" "$f"
    done
fi

# 3. Prune stubs if KEEP is specified
if [ -n "${KEEP_MODULES}" ]; then
    echo "Pruning stubs (keeping only: ${KEEP_MODULES})..."
    ALL_STUBS=("iam" "tenant" "notification" "audit")
    
    # Convert KEEP_MODULES to array
    IFS=',' read -r -a KEEP_ARRAY <<< "${KEEP_MODULES}"
    
    for stub in "${ALL_STUBS[@]}"; do
        keep=false
        for k in "${KEEP_ARRAY[@]}"; do
            if [ "$stub" == "$k" ]; then
                keep=true
                break
            fi
        done
        
        if [ "$keep" = false ]; then
            echo "Removing unused stub: modules/${stub}..."
            rm -rf "${ROOT_DIR}/modules/${stub}"
            
            # Remove from cmd/monolith/modules.go
            if [ -f "${ROOT_DIR}/cmd/monolith/modules.go" ]; then
                sed_i "\|/modules/${stub}\"|d" "${ROOT_DIR}/cmd/monolith/modules.go"
                sed_i "\|${stub}\.New(pool, bus),|d" "${ROOT_DIR}/cmd/monolith/modules.go"
            fi
        fi
    done
fi

echo "Cleaning up dependencies and formatting code..."
go mod tidy
gofmt -w .

echo "Initialization complete!"
echo "Current module: $(grep -E '^module ' go.mod | awk '{print $2}')"
