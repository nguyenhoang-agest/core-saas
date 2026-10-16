 #!/usr/bin/env bash
set -euo pipefail

# ==============================================================================
# Generate a new Goose SQL migration file for a specific module
# Usage: ./scripts/new-migration.sh <module_name> <migration_name>
# Example: ./scripts/new-migration.sh tenant add_custom_domain
# ==============================================================================

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(cd "${SCRIPT_DIR}/.." && pwd)"

if [ $# -lt 2 ]; then
    echo "Usage: $0 <module_name> <migration_name>"
    echo "Example: $0 tenant add_custom_domain"
    exit 1
fi

MODULE_NAME="$1"
MIGRATION_NAME="$2"

MODULE_DIR="${ROOT_DIR}/modules/${MODULE_NAME}"
MIGRATIONS_DIR="${MODULE_DIR}/migrations"

if [ ! -d "${MODULE_DIR}" ]; then
    echo "Error: module directory 'modules/${MODULE_NAME}' does not exist."
    echo "Available modules: $(ls -d "${ROOT_DIR}/modules"/*/ | xargs -n 1 basename | grep -v '_template' | tr '\n' ' ')"
    exit 1
fi

mkdir -p "${MIGRATIONS_DIR}"

# Sanitize migration name: convert spaces/hyphens to underscores, lowercase
MIGRATION_NAME=$(echo "${MIGRATION_NAME}" | tr ' -' '_' | tr '[:upper:]' '[:lower:]' | tr -cd 'a-z0-9_')

if [ -z "${MIGRATION_NAME}" ]; then
    echo "Error: invalid migration name."
    exit 1
fi

# Find next sequential number
LAST_NUM=$(ls "${MIGRATIONS_DIR}"/*.sql 2>/dev/null | xargs -n 1 basename 2>/dev/null | grep -E '^[0-9]{5}_' | sort -n | tail -n 1 | cut -d'_' -f1 || echo "00000")

if [ -z "${LAST_NUM}" ] || [ "${LAST_NUM}" = "00000" ]; then
    NEXT_NUM=1
else
    # Remove leading zeros for arithmetic
    NEXT_NUM=$((10#${LAST_NUM} + 1))
fi

FILE_PREFIX=$(printf "%05d" "${NEXT_NUM}")
TARGET_FILE="${MIGRATIONS_DIR}/${FILE_PREFIX}_${MIGRATION_NAME}.sql"

cat <<EOF > "${TARGET_FILE}"
-- +goose Up
-- +goose StatementBegin
-- Write your schema changes here
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Write rollback statements here
-- +goose StatementEnd
EOF

echo "✓ Created new migration: modules/${MODULE_NAME}/migrations/${FILE_PREFIX}_${MIGRATION_NAME}.sql"
