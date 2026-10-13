# browser-tilion

A composable Docker setup for the **Tilion** browser. It wraps the
[Tilion (Fortress)](https://fortress.tilion.com/) stealth-Chromium engine in a
single image, runs it with `docker-compose`, and exposes the engine's raw CDP
endpoint on `:9222` so any CDP-speaking tool — Playwright, Puppeteer, go-rod, the
Python `tilion` facade — can drive it. A `Makefile` builds the image and manages
the container; the `examples/` directory shows the CDP endpoint in use.

## Why the Tilion browser

Websites detect automation by reading the browser's fingerprint. A stock or
headless Chromium gets flagged — 403s, Cloudflare "Just a moment", DataDome or
PerimeterX walls — before the page even renders. Tilion is built on **Fortress**,
a recompiled Chromium fork that corrects that fingerprint inside Chromium's C++,
with no JavaScript patch layer for a page to catch. The browser reads as an
ordinary Chrome install, so your scraper or agent gets the real page instead of
the interstitial.

## Key features

- **Engine-level stealth.** Canvas, WebGL, WebGPU, audio, fonts, navigator and
  ~30 more surfaces are corrected in C++, not injected via JS, so they hold in
  every frame and worker and survive `toString` / realm re-acquisition.
- **Drop-in over CDP.** Point Playwright, Puppeteer, browser-use or go-rod at
  `http://localhost:9222`; the rest of your code is unchanged. No chromedriver.
- **Python facade.** The `tilion` package boots and drives the engine with an
  async API: `fetch`, `search`, `crawl`, `screenshot`.
- **Agent-ready.** An MCP server exposes 29 tools (`fetch_protected_page`,
  `extract_page`, `crawl_site`, `recon_site_apis`, `search_web`,
  `run_browser_task`, `get_stealth_cdp_endpoint`, …) to Claude, Cursor, Codex and
  other MCP clients. Extras: `pip install "tilion[mcp]"`, `"tilion[vision]"`.
- **Egress control.** `TILION_PROXY` routes traffic through a residential/mobile
  proxy; `TILION_REGION` aligns timezone and locale to the IP.

## How it differs from stock Chromium

Bot detectors look at three layers. Fortress fixes the one that matters, and is
built to keep the other two clean:

| Layer | What leaks | Stock Chromium | Tilion (Fortress) |
| --- | --- | --- | --- |
| Driver / binary artifacts | `cdc_` vars, WebDriver surface | present with chromedriver | drive raw CDP, no chromedriver |
| CDP side-effects | `Runtime.enable` leak, init-script footprints | leaked | no leak (raw CDP, verified) |
| Fingerprint surface | canvas, WebGL, audio, fonts, navigator | flags as headless bot | **corrected in C++** |

JS stealth plugins (`puppeteer-stealth`, `undetected-chromedriver`) patch *after*
the page can inspect the browser and self-reveal. Tilion moves the spoof into the
engine, where the page finds native code.

## What Tilion can do

- Fetch pages behind Cloudflare / DataDome / PerimeterX / Akamai / 403 / CAPTCHA.
- Web search through a real browser, without a SERP API key.
- Crawl a whole site (auto-handling SPA/JS) and reverse-engineer its private
  JSON/XHR API.
- Screenshot pages, extract structured data, and convert PDF/DOCX/XLSX to text.
- Drive a page: click, fill, press, wait, evaluate JS.
- Run multi-step flows: login, paginate, infinite-scroll, checkout.
- Persist and restore authenticated sessions (cookies, save/load profile) and
  download files.
- Hand out a stealth CDP endpoint for your own Playwright/Puppeteer/browser-use
  code.

## What it cannot do

- **It is not a blanket bypass.** Anti-bot systems also weigh IP reputation and
  behavior. On the hardest targets (Cloudflare, DataDome, Kasada) a datacenter IP
  still gets blocked — the fix is residential/mobile egress via `TILION_PROXY`,
  not a smarter fingerprint.
- Passing client-side checks (CreepJS, Sannysoft, BrowserScan) is necessary, not
  sufficient — results move as detectors change.
- It is not driven through Selenium/chromedriver; drive it over CDP instead.
- CAPTCHA solving requires a `CAPTCHA_API_KEY`.
- Behavior, rate and request patterns are still detectable regardless of the
  fingerprint.
- The CDP endpoint is published **without authentication** (see Security).

## Quick start

```bash
make up          # build the image and start the container
curl -s http://localhost:9222/json/version   # engine is reachable
make down        # stop and remove
```

Then drive it from any CDP client, e.g. point Playwright at
`connect_over_cdp("http://localhost:9222")`, or use the bundled `examples/`.

## Security

CDP has no authentication: anyone who can reach the port has full control of the
browser. This project publishes `TILION_PORT` (default `9222`) on `0.0.0.0`.
Keep it closed at the firewall/network level and never expose it to the public
internet. Secrets (`TILION_API_KEY`, LLM keys, `CAPTCHA_API_KEY`) belong in a
git-ignored `.env`, not in the image.
