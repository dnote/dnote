#!/bin/bash
set -e

# Usage: ./generate-changelog.sh <component> <current-tag> <previous-tag>
# Example: ./generate-changelog.sh cli cli-v0.15.2 cli-v0.15.1

COMPONENT=$1
CURRENT_TAG=$2
PREV_TAG=$3

if [ -z "$COMPONENT" ] || [ -z "$CURRENT_TAG" ]; then
  echo "Usage: $0 <component> <current-tag> [previous-tag]"
  echo "Example: $0 cli cli-v0.15.2 cli-v0.15.1"
  exit 1
fi

# Validate that tags match the component
if [[ ! "$CURRENT_TAG" =~ ^${COMPONENT}- ]]; then
  echo "Error: Current tag '$CURRENT_TAG' doesn't match component '$COMPONENT'"
  echo "Expected tag to start with '${COMPONENT}-'"
  exit 1
fi

if [ -n "$PREV_TAG" ] && [[ ! "$PREV_TAG" =~ ^${COMPONENT}- ]]; then
  echo "Error: Previous tag '$PREV_TAG' doesn't match component '$COMPONENT'"
  echo "Expected tag to start with '${COMPONENT}-'"
  exit 1
fi

# Define paths for each component
# Shared paths that apply to both components
SHARED_PATHS="pkg/dirs/"

if [ "$COMPONENT" == "cli" ]; then
  FILTER_PATHS="pkg/cli/ cmd/cli/ $SHARED_PATHS"
elif [ "$COMPONENT" == "server" ]; then
  FILTER_PATHS="pkg/server/ host/ $SHARED_PATHS"
else
  echo "Unknown component: $COMPONENT"
  echo "Valid components: cli, server"
  exit 1
fi

# Determine commit range
if [ -z "$PREV_TAG" ]; then
  echo "Error: No previous tag specified"
  exit 1
fi

RANGE="$PREV_TAG..$CURRENT_TAG"

# Get all commits that touched the relevant paths
# Warnings go to stderr (visible in logs), commits go to stdout (captured for file)
COMMITS=$(
  for path in $FILTER_PATHS; do
    git log --oneline --no-merges --pretty=format:"- %s%n" "$RANGE" -- "$path" 2>&2
  done | sort -u | grep -v "^$"
)

echo "## What's Changed"
echo ""
echo "$COMMITS"
