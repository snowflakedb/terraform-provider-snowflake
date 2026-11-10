#!/usr/bin/env bash
set -euo pipefail

if [[ $# -ne 1 ]]; then
  echo "Usage: $0 <version>" >&2
  echo "  e.g. $0 2.22.1" >&2
  exit 1
fi

VERSION="${1#v}"   # accept 2.22.1 or v2.22.1
REPO=snowflakedb/terraform-provider-snowflake

RELEASE_OWNER="$(git config get user.email)" || { echo "git user.email is not set" >&2; exit 1; }

SHA="$(gh api "repos/${REPO}/commits/refs/tags/v${VERSION}" --jq .sha)" \
  || { echo "Tag v${VERSION} not found in ${REPO}" >&2; exit 1; }

CONTROLLER=https://es-ci-legacy-mainvalidation-001.jenkinsdev1.us-west-2.aws-dev.app.snowflake.com
SNYK=https://app.snyk.io/org/snowflakedb-sca-scanning-private-repo/project/d572d366-ebce-40e9-82df-787c845caf53
RELEASE_URL="https://github.com/${REPO}/releases/tag/v${VERSION}"

enc() { jq -rn --arg v "$1" '$v|@uri'; }

PARAMS="DRIVER_CLIENT_CONNECTOR_TYPE=terraform-provider-snowflake"
PARAMS+="&RELEASE_VERSION=$(enc "$VERSION")"
PARAMS+="&COMMIT_SHA=$(enc "$SHA")"
PARAMS+="&REGRESSION_LINKS=N%2FA"
PARAMS+="&PERFORMANCE_LINKS=N%2FA"
PARAMS+="&SECURITY_VULNERABILITY=$(enc "No vulnerabilities.
$SNYK")"
PARAMS+="&EPIC_JIRA=SNOW-2236295"
PARAMS+="&CLIENT_RELEASE_NOTES=$(enc "$RELEASE_URL")"
PARAMS+="&RELEASE_SIGN_OFF=$(enc "$RELEASE_OWNER")"
PARAMS+="&RELEASE_NOTES_OWNER=$(enc "$RELEASE_OWNER")"
PARAMS+="&RELEASE_OWNER=$(enc "$RELEASE_OWNER")"
PARAMS+="&START_RELEASE=false"

sf jenkins submit \
  -c "$CONTROLLER" \
  -j RequestClientRelease \
  -p "$PARAMS"
