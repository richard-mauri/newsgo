# NewsGo v0.2.4 — Public reader URL and dedicated launcher

NewsGo now defaults to `https://scrapergo.themauris.org` for headline links. Override for development with `-scrapergo http://127.0.0.1:8083`. Named tab reuse (`newsgo-scrapergo`) is preserved.

Build the native Intel binary from `~/projects/newsgo`:

```sh
mkdir -p dist .cloudflared
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o dist/newsgo-darwin-amd64 ./cmd/newsgo
```

Place the **NewsGo tunnel token** in `.cloudflared/token` (not the ScraperGo token). Configure the public hostname `newsgo.themauris.org` to route to `http://127.0.0.1:8084` in Cloudflare. Then:

```sh
chmod 600 .cloudflared/token
chmod +x newsgo.sh
./newsgo.sh
curl -I http://127.0.0.1:8084/
```

Do not commit `.cloudflared/token` or local logs; add them to your local `.gitignore` if not already excluded. `newsgo.sh` intentionally does not modify another project's files or tunnels. The Go web server binds to loopback only. FeedCheck and personal subscriptions remain unchanged.
