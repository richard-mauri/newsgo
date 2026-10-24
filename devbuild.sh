#!/usr/bin/env bash
# NewsGo: build standalone macOS Intel/Apple Silicon binaries; optionally publish a GitHub Release.
set -Eeuo pipefail

SCRIPT_DIR="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"
REPO="richard-mauri/newsgo"
OUT_DIR="$SCRIPT_DIR/dist"
VERSION=""

usage() {
  cat <<'HELP'
Usage:
  ./devbuild.sh                  Build both macOS binaries and SHA256SUMS
  ./devbuild.sh --release vX.Y.Z Build and publish a new GitHub Release
  ./devbuild.sh --help           Display this help

Requires Go for building. Publishing additionally requires an authenticated
GitHub CLI (gh), a clean Git working tree, and the current HEAD pushed to its
upstream branch. A release with the requested tag must not already exist.

Release assets: dist/newsgo-darwin-amd64,
                dist/newsgo-darwin-arm64, dist/SHA256SUMS
HELP
}

case "${1:-}" in
  "") [[ $# -eq 0 ]] || { usage >&2; exit 2; } ;;
  -h|--help) [[ $# -eq 1 ]] || { usage >&2; exit 2; }; usage; exit 0 ;;
  --release)
    [[ $# -eq 2 ]] || { echo 'Error: expected --release vX.Y.Z' >&2; usage >&2; exit 2; }
    VERSION="$2"
    [[ "$VERSION" =~ ^v[0-9]+\.[0-9]+\.[0-9]+$ ]] || {
      echo 'Error: release version must look like v0.2.4' >&2; exit 2;
    } ;;
  *) echo "Error: unknown option: $1" >&2; usage >&2; exit 2 ;;
esac

command -v go >/dev/null 2>&1 || { echo 'Error: Go is not installed.' >&2; exit 1; }
[[ -f go.mod ]] || { echo "Error: missing go.mod in $SCRIPT_DIR" >&2; exit 1; }
[[ -f cmd/newsgo/main.go ]] || { echo 'Error: missing cmd/newsgo/main.go' >&2; exit 1; }

if [[ -n "$VERSION" ]]; then
  for cmd in gh git; do
    command -v "$cmd" >/dev/null 2>&1 || { echo "Error: $cmd is required to publish a release." >&2; exit 1; }
  done
  gh auth status >/dev/null 2>&1 || { echo 'Error: authenticate with gh auth login first.' >&2; exit 1; }
  [[ -z "$(git status --porcelain)" ]] || { echo 'Error: commit or stash changes before publishing a release.' >&2; exit 1; }
  branch="$(git symbolic-ref --quiet --short HEAD || true)"
  [[ "$branch" == main ]] || { echo 'Error: publish releases from the main branch.' >&2; exit 1; }
  upstream="$(git rev-parse --abbrev-ref --symbolic-full-name '@{upstream}' 2>/dev/null || true)"
  [[ -n "$upstream" ]] || { echo 'Error: no upstream configured. Run git push -u origin main.' >&2; exit 1; }
  [[ "$(git rev-parse HEAD)" == "$(git rev-parse '@{upstream}')" ]] || {
    echo 'Error: local HEAD differs from its upstream. Commit and git push first.' >&2; exit 1;
  }
  # Also confirm the remote branch matches HEAD (a local tracking ref can be stale).
  remote_head="$(git ls-remote origin refs/heads/main | awk '{print $1}')"
  [[ -n "$remote_head" && "$(git rev-parse HEAD)" == "$remote_head" ]] || {
    echo 'Error: HEAD does not match origin/main on GitHub. Push changes first.' >&2; exit 1;
  }
  if gh release view "$VERSION" --repo "$REPO" >/dev/null 2>&1; then
    echo "Error: GitHub Release $VERSION already exists. Choose a new version." >&2; exit 1
  fi
  if git ls-remote --exit-code --tags origin "refs/tags/$VERSION" >/dev/null 2>&1; then
    echo "Error: Git tag $VERSION already exists on origin. Choose a new version." >&2; exit 1
  fi
fi

counter_path="$(git rev-parse --git-path newsgo-build-counter)"
[[ "$counter_path" = /* ]] || counter_path="$SCRIPT_DIR/$counter_path"
previous=0
if [[ -f "$counter_path" ]]; then
  IFS= read -r saved < "$counter_path" || true
  [[ "$saved" =~ ^[0-9]{1,6}$ ]] || { echo "Invalid local build counter" >&2; exit 1; }
  previous=$((10#$saved))
fi
(( previous < 999999 )) || { echo "Build counter exhausted" >&2; exit 1; }
next_build="$(printf '%06d' "$((previous+1))")"
echo "NewsGo candidate Build $next_build"
mkdir -p "$OUT_DIR"
printf 'NewsGo build: %s\n' "$(go version)"
printf 'Running tests...\n'
go test ./...

for arch in amd64 arm64; do
  output="$OUT_DIR/newsgo-darwin-$arch"
  printf 'Building darwin/%s...\n' "$arch"
  CGO_ENABLED=0 GOOS=darwin GOARCH="$arch" \
    go build -trimpath -buildvcs=false -ldflags="-s -w -X github.com/richard-mauri/newsgo/internal/dashboard.BuildNumber=$next_build" -o "$output" ./cmd/newsgo
  chmod 755 "$output"
  command -v file >/dev/null 2>&1 && file "$output"
done

(
  cd "$OUT_DIR"
  shasum -a 256 newsgo-darwin-amd64 newsgo-darwin-arm64 > SHA256SUMS
)
printf '%s\n' "$next_build" > "$counter_path"
printf '\nBuilt binaries and checksums in %s (Build %s)\n' "$OUT_DIR" "$next_build"
ls -lh "$OUT_DIR/newsgo-darwin-amd64" "$OUT_DIR/newsgo-darwin-arm64" "$OUT_DIR/SHA256SUMS"

if [[ -n "$VERSION" ]]; then
  printf '\nPublishing GitHub Release %s...\n' "$VERSION"
  gh release create "$VERSION" \
    "$OUT_DIR/newsgo-darwin-amd64" \
    "$OUT_DIR/newsgo-darwin-arm64" \
    "$OUT_DIR/SHA256SUMS" \
    --repo "$REPO" --target main \
    --title "NewsGo $VERSION" \
    --generate-notes
  printf '\nRelease published: https://github.com/%s/releases/tag/%s\n' "$REPO" "$VERSION"
else
  printf '\nTo publish after committing and pushing changes:\n  ./devbuild.sh --release vX.Y.Z\n'
fi
