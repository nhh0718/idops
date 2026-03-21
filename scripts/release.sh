#!/usr/bin/env bash
# Quick release script - tag + push to trigger GitHub Actions release
set -euo pipefail

VERSION="${1:-}"

if [[ -z "$VERSION" ]]; then
    # Show current tags
    echo "Current tags:"
    git tag --sort=-version:refname | head -5 2>/dev/null || echo "  (none)"
    echo ""
    read -rp "Enter new version (e.g. v0.1.0): " VERSION
fi

# Validate version format
if [[ ! "$VERSION" =~ ^v[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
    echo "Error: Version must match vX.Y.Z (e.g. v0.1.0)"
    exit 1
fi

# Check if tag already exists
if git rev-parse "$VERSION" >/dev/null 2>&1; then
    echo "Error: Tag $VERSION already exists"
    exit 1
fi

# Check for uncommitted changes
if [[ -n "$(git status --porcelain)" ]]; then
    echo "Warning: You have uncommitted changes."
    read -rp "Continue anyway? [y/N]: " CONFIRM
    [[ "$CONFIRM" =~ ^[yY]$ ]] || exit 1
fi

echo ""
echo "=== Release $VERSION ==="
echo ""

# Ensure we're on main
BRANCH="$(git branch --show-current)"
echo "  Branch:  $BRANCH"
echo "  Version: $VERSION"
echo ""

read -rp "Create tag $VERSION and push? [y/N]: " CONFIRM
[[ "$CONFIRM" =~ ^[yY]$ ]] || exit 1

# Create tag
git tag -a "$VERSION" -m "Release $VERSION"
echo "  Created tag $VERSION"

# Push tag (triggers GitHub Actions release workflow)
git push origin "$VERSION"
echo "  Pushed tag to origin"

echo ""
echo "Done! GitHub Actions will build and release automatically."
echo "Check: https://github.com/$(git remote get-url origin | sed 's/.*github.com[:/]\(.*\)\.git/\1/')/actions"
