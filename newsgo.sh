#!/bin/bash
# Start NewsGo (Intel/amd64) and its own Cloudflare Tunnel.
# ScraperGo and all other services/tunnels remain untouched.
# Requires newsgo/.cloudflared/token (do not commit this file).

HOME=/Users/richardmauri
APP_DIR="$HOME/projects/newsgo"
APP="$APP_DIR/dist/newsgo-darwin-amd64"
TOKEN="$APP_DIR/.cloudflared/token"
CLOUDFLARED="$HOME/Applications/cloudflared"

cd "$APP_DIR" || exit 1
for file in "$APP" "$TOKEN" "$CLOUDFLARED"; do
    if [[ ! -f "$file" ]]; then
        echo "Missing required file: $file" >&2
        exit 1
    fi
done
if [[ ! -x "$APP" || ! -x "$CLOUDFLARED" ]]; then
    echo "NewsGo binary and cloudflared must be executable." >&2
    exit 1
fi

# Verify the *exact* executable before stopping any existing NewsGo process.
# An old dist binary without -version must not silently replace a working app.
echo "NewsGo executable: $APP"
if ! version_output="$("$APP" -version 2>&1)"; then
    printf 'Error: NewsGo version check failed: %s\n' "$version_output" >&2
    exit 1
fi
if [[ ! "$version_output" =~ ^NewsGo\ v[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
    printf 'Error: unexpected NewsGo executable version: %s\n' "$version_output" >&2
    exit 1
fi
printf 'NewsGo executable version: %s\n' "$version_output"

# Stop only the NewsGo instance launched from this project's dist directory.
ps -axo pid=,command= | while read -r pid command; do
    if [[ "$command" == "$APP"* ]]; then
        kill "$pid" 2>/dev/null || true
    fi
done

# Stop only NewsGo's Cloudflare Tunnel (identified by NewsGo's token path).
ps -axo pid=,command= | while read -r pid command; do
    if [[ "$command" == *"cloudflared tunnel run"* &&
          "$command" == *"$TOKEN"* ]]; then
        kill "$pid" 2>/dev/null || true
    fi
done

# Bind locally; the public hostname is routed through NewsGo's own tunnel.
nohup "$APP" -listen 127.0.0.1:8084 \
    -scrapergo https://scrapergo.themauris.org \
    > "$APP_DIR/newsgo.log" 2>&1 < /dev/null &

nohup "$CLOUDFLARED" tunnel run \
    --token-file "$TOKEN" \
    > "$APP_DIR/.cloudflared/tunnel.log" 2>&1 < /dev/null &

disown -a

echo "Started NewsGo on 127.0.0.1:8084 and its dedicated Cloudflare Tunnel."
