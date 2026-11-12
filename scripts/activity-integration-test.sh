#!/usr/bin/env bash
# Isolated NewsGo -> Homepage activity integration smoke test.
# No launchers, no Cloudflare tunnels, no production data or process changes.
set -Eeuo pipefail

usage() {
  echo 'Usage: TEST_ROOT=/path/to/homepage-activity.XXXXXX bash scripts/activity-integration-test.sh' >&2
  echo 'TEST_ROOT must be the EXISTING isolated Homepage test root (containing homepage and data/).' >&2
}

if [[ -z "${TEST_ROOT:-}" ]]; then usage; exit 2; fi
case "$TEST_ROOT" in
  */homepage-activity.*) ;;
  *) echo 'Refusing a path not named homepage-activity.*' >&2; exit 2 ;;
esac
if [[ ! -d "$TEST_ROOT/data" || ! -x "$TEST_ROOT/homepage" || -L "$TEST_ROOT/data" ]]; then
  echo 'Expected existing isolated Homepage binary and non-symlink data directory.' >&2
  exit 2
fi
if [[ ! -f "$TEST_ROOT/data/auth/approved-tokens.json" ]]; then
  echo 'Missing disposable test admin. See Homepage scripts/prepare-activity-ui-test.sh.' >&2
  exit 2
fi
for binary in curl lsof openssl go; do
  command -v "$binary" >/dev/null || { echo "Required tool unavailable: $binary" >&2; exit 1; }
done
for port in 8788 8085; do
  if lsof -nP -iTCP:"$port" -sTCP:LISTEN | grep -q LISTEN; then
    echo "Port $port is occupied. Stop only the previous ISOLATED test server there; do not stop production services." >&2
    exit 2
  fi
done
repo="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
test_dir="$(mktemp -d "${TMPDIR:-/tmp}/newsgo-activity.XXXXXX")"
key="$(openssl rand -hex 32)"
home_pid=''
news_pid=''

cleanup() {
  local pid
  for pid in "$news_pid" "$home_pid"; do
    if [[ -n "$pid" ]]; then
      kill "$pid" 2>/dev/null || true
      wait "$pid" 2>/dev/null || true
    fi
  done
  echo "Isolated servers stopped. Diagnostic logs remain in $test_dir"
}
trap cleanup EXIT

echo 'Compiling isolated NewsGo candidate without touching the installed executable.'
(cd "$repo" && go test -count=1 ./... && go build -o "$test_dir/newsgo" ./cmd/newsgo)
mkdir -p "$test_dir/settings"
echo 'Starting isolated Homepage on 127.0.0.1:8788, retaining its existing test history.'
(
  cd "$TEST_ROOT/src"
  SUITE_DATA_DIR="$TEST_ROOT/data" HOMEPAGE_ACTIVITY_KEY_NEWSGO="$key" \
    "$TEST_ROOT/homepage" > "$test_dir/homepage.log" 2>&1
) &
home_pid=$!
echo 'Starting isolated NewsGo on 127.0.0.1:8085.'
(
  cd "$repo"
  NEWSGO_ACTIVITY_URL="http://127.0.0.1:8788/api/admin/activity/ingest" \
  HOMEPAGE_ACTIVITY_KEY_NEWSGO="$key" \
    "$test_dir/newsgo" \
      -listen '127.0.0.1:8085' \
      -feeds "$repo/feeds.json" \
      -settings-dir "$test_dir/settings" \
      -scrapergo 'http://127.0.0.1:1' \
    > "$test_dir/newsgo.log" 2>&1
) &
news_pid=$!

sleep 2
if ! curl --fail --silent --show-error http://127.0.0.1:8788/healthz >/dev/null; then
  echo "Homepage test health failed. Inspect $test_dir/homepage.log" >&2; exit 1
fi
if ! curl --fail --silent --show-error http://127.0.0.1:8085/health >/dev/null; then
  echo "NewsGo test health failed. Inspect $test_dir/newsgo.log" >&2; exit 1
fi
cat <<'MESSAGE'

Both ISOLATED services are responding.
Open Safari:
  NewsGo:            http://127.0.0.1:8085/
  Homepage admin:    http://127.0.0.1:8788/suite-admin

Sign in at the Homepage admin page with the EXISTING disposable test token,
then keep the NewsGo tab visible for at least 75 seconds.
Click Refresh activity. NewsGo should change from Unknown to Reporting,
show a nonzero Active in 2 minutes count, and add hourly history samples.

This test deliberately directs NewsGo's ScraperGo reader to unused port 1.
Do not test real feed/reader extraction with this harness.

Use Ctrl-C in this terminal to terminate ONLY the two isolated test servers.
No production service, launch script or Cloudflare tunnel is touched.
MESSAGE
while kill -0 "$home_pid" 2>/dev/null && kill -0 "$news_pid" 2>/dev/null; do
  sleep 2
done
echo 'One isolated test server exited; inspect the log paths shown above.' >&2
exit 1
