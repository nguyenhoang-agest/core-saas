#!/usr/bin/env bash
set -euo pipefail

# Usage: 
#   ./scripts/new-module.sh <module_name> [--dirs-only]
#   DIRS_ONLY=true ./scripts/new-module.sh <module_name>

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(cd "${SCRIPT_DIR}/.." && pwd)"

if [ $# -lt 1 ]; then
    echo "Usage: $0 <module_name> [--dirs-only]"
    echo "Example: $0 billing"
    echo "Example (structure only): $0 billing --dirs-only"
    exit 1
fi

MODULE_NAME="$1"

# Parse --dirs-only / --bare option
DIRS_ONLY="${DIRS_ONLY:-false}"
for arg in "$@"; do
    case "$arg" in
        --dirs-only|--bare|--skeleton|dirs-only|bare|skeleton|true|1)
            DIRS_ONLY=true
            ;;
    esac
done

# Validate module name
if [[ ! "${MODULE_NAME}" =~ ^[a-z][a-z0-9_]*$ ]]; then
    echo "Error: module name must start with a lowercase letter and contain only lowercase alphanumeric characters and underscores."
    exit 1
fi

TARGET_DIR="${ROOT_DIR}/modules/${MODULE_NAME}"
TEMPLATE_DIR="${ROOT_DIR}/modules/_template"
REGISTRY_FILE="${ROOT_DIR}/cmd/monolith/modules.go"

if [ -d "${TARGET_DIR}" ]; then
    echo "Error: module directory 'modules/${MODULE_NAME}' already exists."
    exit 1
fi

# Detect current Go module path from go.mod
MOD_PATH=$(grep -E '^module ' "${ROOT_DIR}/go.mod" | awk '{print $2}')
if [ -z "${MOD_PATH}" ]; then
    echo "Error: unable to determine Go module path from go.mod"
    exit 1
fi

# Helper functions for case formatting
to_pascal() {
    local input="$1"
    echo "$input" | awk -F'_' '{for(i=1;i<=NF;i++) printf toupper(substr($i,1,1)) substr($i,2); print ""}'
}

to_plural() {
    local input="$1"
    if [[ "$input" == *s ]]; then
        echo "${input}"
    elif [[ "$input" == *y && "$input" != *ay && "$input" != *ey && "$input" != *oy && "$input" != *uy ]]; then
        echo "${input%y}ies"
    else
        echo "${input}s"
    fi
}

ENTITY_NAME=$(to_pascal "${MODULE_NAME}")
PLURAL_NAME=$(to_plural "${MODULE_NAME}")
PLURAL_TITLE=$(to_pascal "${PLURAL_NAME}")

# ==============================================================================
# MODE 1: Directory Structure Only (No Code Generation)
# ==============================================================================
if [ "${DIRS_ONLY}" = "true" ]; then
    echo "Scaffolding directory structure for '${MODULE_NAME}' (structure only, no code)..."

    mkdir -p "${TARGET_DIR}/api"
    mkdir -p "${TARGET_DIR}/internal/dto"
    mkdir -p "${TARGET_DIR}/internal/errcode"
    mkdir -p "${TARGET_DIR}/internal/handler"
    mkdir -p "${TARGET_DIR}/internal/model"
    mkdir -p "${TARGET_DIR}/internal/repository"
    mkdir -p "${TARGET_DIR}/internal/service"
    mkdir -p "${TARGET_DIR}/migrations"

    touch "${TARGET_DIR}/api/.gitkeep"
    touch "${TARGET_DIR}/internal/dto/.gitkeep"
    touch "${TARGET_DIR}/internal/errcode/.gitkeep"
    touch "${TARGET_DIR}/internal/handler/.gitkeep"
    touch "${TARGET_DIR}/internal/model/.gitkeep"
    touch "${TARGET_DIR}/internal/repository/.gitkeep"
    touch "${TARGET_DIR}/internal/service/.gitkeep"
    touch "${TARGET_DIR}/migrations/.gitkeep"

    cat <<EOF > "${TARGET_DIR}/README.md"
# ${ENTITY_NAME} Module

## Structure
- \`api/\`: Public contracts, DTOs, and event topics
- \`internal/dto/\`: Request/Response data transfer objects
- \`internal/errcode/\`: Standardized machine-readable error codes
- \`internal/handler/\`: HTTP handlers and Fiber routing
- \`internal/model/\`: Domain entities and database models
- \`internal/repository/\`: Data access layer
- \`internal/service/\`: Business domain use cases
- \`migrations/\`: Goose SQL migrations

## Registration
To register this module into the monolith:
1. Implement \`kernel.Module\` in \`modules/${MODULE_NAME}/module.go\`
2. Register \`${MODULE_NAME}.New(pool, bus)\` in \`cmd/monolith/modules.go\`
See [docs/ADDING_A_MODULE.md](../../docs/ADDING_A_MODULE.md).
EOF

    echo "✓ Successfully created directory structure for 'modules/${MODULE_NAME}'."
    echo ""
    echo "Directory layout:"
    echo "  modules/${MODULE_NAME}/"
    echo "  ├── api/"
    echo "  ├── internal/"
    echo "  │   ├── dto/"
    echo "  │   ├── errcode/"
    echo "  │   ├── handler/"
    echo "  │   ├── model/"
    echo "  │   ├── repository/"
    echo "  │   └── service/"
    echo "  ├── migrations/"
    echo "  └── README.md"
    echo ""
    echo "Next Steps:"
    echo "  1. Add your code inside modules/${MODULE_NAME}/"
    echo "  2. Implement kernel.Module in modules/${MODULE_NAME}/module.go"
    echo "  3. Register ${MODULE_NAME}.New(pool, bus) in cmd/monolith/modules.go"
    echo "  (See docs/ADDING_A_MODULE.md for registration details)"
    exit 0
fi

# ==============================================================================
# MODE 2: Full Code Generation (from _template)
# ==============================================================================
if [ ! -d "${TEMPLATE_DIR}" ]; then
    echo "Error: template directory 'modules/_template' not found."
    exit 1
fi

echo "Scaffolding module '${MODULE_NAME}' based on Tenant pattern..."
echo "  Entity Name : ${ENTITY_NAME}"
echo "  Plural Name : ${PLURAL_NAME}"
echo "  Route Group : /${PLURAL_NAME}"
echo "  Table Name  : ${PLURAL_NAME}"

cp -r "${TEMPLATE_DIR}" "${TARGET_DIR}"

# Rename entity model and migration files
if [ -f "${TARGET_DIR}/internal/model/template.go" ]; then
    mv "${TARGET_DIR}/internal/model/template.go" "${TARGET_DIR}/internal/model/${MODULE_NAME}.go"
fi

if [ -f "${TARGET_DIR}/migrations/00001_create_templates_table.sql" ]; then
    mv "${TARGET_DIR}/migrations/00001_create_templates_table.sql" "${TARGET_DIR}/migrations/00001_create_${PLURAL_NAME}_table.sql"
fi

# Use platform-compatible sed
sed_i() {
    if [[ "$OSTYPE" == "darwin"* ]]; then
        sed -i '' "$@"
    else
        sed -i "$@"
    fi
}

# Perform template substitutions across all files in TARGET_DIR
find "${TARGET_DIR}" -type f \( -name "*.go" -o -name "*.sql" -o -name "*.md" \) | while read -r file; do
    # 1. Update module import paths
    sed_i "s|${MOD_PATH}/modules/_template|${MOD_PATH}/modules/${MODULE_NAME}|g" "${file}"
    # 2. Plural titles (e.g. Templates -> Billings)
    sed_i "s|Templates|${PLURAL_TITLE}|g" "${file}"
    # 3. Plural names (e.g. templates -> billings)
    sed_i "s|templates|${PLURAL_NAME}|g" "${file}"
    # 4. Entity names (e.g. Template -> Billing)
    sed_i "s|Template|${ENTITY_NAME}|g" "${file}"
    # 5. Lowercase names (e.g. template -> billing)
    sed_i "s|template|${MODULE_NAME}|g" "${file}"
done

# Register in cmd/monolith/modules.go if it exists
if [ -f "${REGISTRY_FILE}" ]; then
    echo "Registering '${MODULE_NAME}' in cmd/monolith/modules.go..."
    
    # Add import if not present
    if ! grep -q "\"${MOD_PATH}/modules/${MODULE_NAME}\"" "${REGISTRY_FILE}"; then
        sed_i "/import (/a\\
	\"${MOD_PATH}/modules/${MODULE_NAME}\"" "${REGISTRY_FILE}"
    fi
    
    # Add to enabledModules slice if not present
    if ! grep -q "${MODULE_NAME}.New(pool, bus)" "${REGISTRY_FILE}"; then
        sed_i "/return \[\]kernel.Module{/a\\
		${MODULE_NAME}.New(pool, bus)," "${REGISTRY_FILE}"
    fi
fi

# Format code
cd "${ROOT_DIR}"
gofmt -w "${TARGET_DIR}" "${REGISTRY_FILE}" >/dev/null 2>&1 || true

echo "✓ Successfully created module 'modules/${MODULE_NAME}' adhering to Tenant architecture."
echo "Running module tests..."
go test -v "./modules/${MODULE_NAME}/..."

echo ""
echo "All tests passed for 'modules/${MODULE_NAME}'!"
echo "Registered in: cmd/monolith/modules.go"
echo "Ready for development. Run 'make dev-backend' or 'make build' to compile."
