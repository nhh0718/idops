#!/usr/bin/env bash
# Build script for idops - cross-platform build with version injection
set -euo pipefail

VERSION="${1:-$(git describe --tags --always --dirty 2>/dev/null || echo "dev")}"
COMMIT="$(git rev-parse --short HEAD 2>/dev/null || echo "none")"
DATE="$(date -u +"%Y-%m-%dT%H:%M:%SZ")"
BINARY="idops"
OUTPUT_DIR="bin"

LDFLAGS="-s -w -X main.version=${VERSION} -X main.commit=${COMMIT} -X main.date=${DATE}"

echo "=== idops build ==="
echo "  Version: ${VERSION}"
echo "  Commit:  ${COMMIT}"
echo "  Date:    ${DATE}"
echo ""

build() {
    local os=$1 arch=$2
    local ext=""
    [[ "$os" == "windows" ]] && ext=".exe"
    local output="${OUTPUT_DIR}/${BINARY}_${os}_${arch}${ext}"

    echo "  Building ${os}/${arch} -> ${output}"
    GOOS=$os GOARCH=$arch go build -ldflags "${LDFLAGS}" -o "${output}" ./cmd/idops
}

mkdir -p "${OUTPUT_DIR}"

case "${2:-local}" in
    local)
        echo "Building for current platform..."
        go build -ldflags "${LDFLAGS}" -o "${OUTPUT_DIR}/${BINARY}" ./cmd/idops
        echo ""
        echo "Done: ${OUTPUT_DIR}/${BINARY}"
        echo "Run:  ./${OUTPUT_DIR}/${BINARY} --help"
        ;;
    all)
        echo "Cross-compiling for all platforms..."
        build linux amd64
        build linux arm64
        build darwin amd64
        build darwin arm64
        build windows amd64
        build windows arm64
        echo ""
        echo "Done. Binaries in ${OUTPUT_DIR}/"
        ls -lh "${OUTPUT_DIR}/"
        ;;
    *)
        echo "Usage: $0 [version] [local|all]"
        echo "  $0              # build for current platform"
        echo "  $0 v1.0.0 all   # cross-compile all platforms"
        exit 1
        ;;
esac
