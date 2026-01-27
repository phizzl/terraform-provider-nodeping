#!/bin/bash
set -e

VERSION="${1:-0.2.0}"
GITLAB_URL="${GITLAB_URL:-https://gitlab.nxs360.com}"
GITLAB_API_URL="${GITLAB_API_URL:-${GITLAB_URL}/api/v4}"
PROJECT_ID="${PROJECT_ID:-1425}"

OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case "$ARCH" in
  x86_64) ARCH="amd64" ;;
  aarch64|arm64) ARCH="arm64" ;;
  *) echo "Unsupported architecture: $ARCH"; exit 1 ;;
esac

BINARY_NAME="terraform-provider-nodeping_${OS}_${ARCH}"
if [ "$OS" = "windows" ]; then
  BINARY_NAME="${BINARY_NAME}.exe"
fi

INSTALL_DIR="${HOME}/.terraform.d/plugins/gitlab.nxs360.com/devops/nodeping/${VERSION}/${OS}_${ARCH}"

echo "Installing terraform-provider-nodeping v${VERSION} for ${OS}/${ARCH}..."

mkdir -p "$INSTALL_DIR"

if [ -n "$CI_JOB_TOKEN" ]; then
  AUTH_HEADER="JOB-TOKEN: ${CI_JOB_TOKEN}"
elif [ -n "$GITLAB_TOKEN" ]; then
  AUTH_HEADER="PRIVATE-TOKEN: ${GITLAB_TOKEN}"
else
  echo "Error: Set GITLAB_TOKEN or CI_JOB_TOKEN for authentication"
  exit 1
fi

curl -fsSL --header "$AUTH_HEADER" \
  "${GITLAB_API_URL}/projects/${PROJECT_ID}/packages/generic/terraform-provider-nodeping/${VERSION}/${BINARY_NAME}" \
  -o "${INSTALL_DIR}/terraform-provider-nodeping"

chmod +x "${INSTALL_DIR}/terraform-provider-nodeping"

echo "Installed to: ${INSTALL_DIR}/terraform-provider-nodeping"
echo ""
echo "Add to your Terraform configuration:"
echo ""
echo 'terraform {'
echo '  required_providers {'
echo '    nodeping = {'
echo '      source  = "gitlab.nxs360.com/devops/nodeping"'
echo "      version = \"${VERSION}\""
echo '    }'
echo '  }'
echo '}'
