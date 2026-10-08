#!/usr/bin/env bash
set -euo pipefail

# Usage: ./scripts/new-module.sh <module_name>

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(cd "${SCRIPT_DIR}/.." && pwd)"

if [ $# -lt 1 ]; then
    echo "Usage: $0 <module_name>"
    echo "Example: $0 billing"
    exit 1
fi

MODULE_NAME="$1"

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

if [ ! -d "${TEMPLATE_DIR}" ]; then
    echo "Error: template directory 'modules/_template' not found."
    exit 1
fi

# Detect current Go module path from go.mod
MOD_PATH=$(grep -E '^module ' "${ROOT_DIR}/go.mod" | awk '{print $2}')
if [ -z "${MOD_PATH}" ]; then
    echo "Error: unable to determine Go module path from go.mod"
    exit 1
fi

echo "Scaffolding module '${MODULE_NAME}' from _template..."
cp -r "${TEMPLATE_DIR}" "${TARGET_DIR}"

# Rename package names and import paths
# Use platform-compatible sed
sed_i() {
    if [[ "$OSTYPE" == "darwin"* ]]; then
        sed -i '' "$@"
    else
        sed -i "$@"
    fi
}

find "${TARGET_DIR}" -type f -name "*.go" | while read -r file; do
    # Replace imports
    sed_i "s|${MOD_PATH}/modules/_template|${MOD_PATH}/modules/${MODULE_NAME}|g" "${file}"
    # Replace package declaration
    sed_i "s|^package template$|package ${MODULE_NAME}|g" "${file}"
    # Replace templatehttp alias
    sed_i "s|templatehttp|${MODULE_NAME}http|g" "${file}"
    # Replace module name string literal
    sed_i "s|return \"template\"|return \"${MODULE_NAME}\"|g" "${file}"
    # Replace route prefix
    sed_i "s|/template|/${MODULE_NAME}|g" "${file}"
done

# Update README
sed_i "s|Template|${MODULE_NAME}|g" "${TARGET_DIR}/README.md"
sed_i "s|template|${MODULE_NAME}|g" "${TARGET_DIR}/README.md"

# Register in cmd/monolith/modules.go if it exists
if [ -f "${REGISTRY_FILE}" ]; then
    echo "Registering '${MODULE_NAME}' in cmd/monolith/modules.go..."
    
    # Add import
    sed_i "/import (/a\\
	\"${MOD_PATH}/modules/${MODULE_NAME}\"" "${REGISTRY_FILE}"
    
    # Add to enabledModules slice
    sed_i "/return \[\]kernel.Module{/a\\
		${MODULE_NAME}.New(pool, bus)," "${REGISTRY_FILE}"
fi

# Format code
cd "${ROOT_DIR}"
gofmt -w "${TARGET_DIR}" "${REGISTRY_FILE}" >/dev/null 2>&1 || true

echo "Successfully created module 'modules/${MODULE_NAME}'."
echo "Registered in: cmd/monolith/modules.go"
echo "Run 'make build' to verify."
