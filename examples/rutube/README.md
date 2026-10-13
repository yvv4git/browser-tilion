# rutube

Reproduces the crawler failure on `https://rutube.ru` and shows where the
page hangs. The crawler (Playwright over the same CDP endpoint) fails on this
site with `Wait page stable: context deadline exceeded` and then
`collect dynamic html: operation timed out`, while a normal browser opens it
fine. This example runs the same sequence the crawler uses — navigate, wait
for load, wait for the page to become stable, collect HTML — and prints
diagnostics at every step:

- a timeline of `document.readyState` right after navigation
- whether the `load` event fires within the timeout
- whether the page becomes stable (load + network idle + DOM quiet) within
  the timeout — this is exactly what the crawler's "Wait page stable" does
- network counters: requests sent / responses / finished / failed, and the
  list of requests that started but never finished (the usual reason "wait
  stable" never returns: a streaming video or a long-lived XHR keeps the
  network busy forever)
- page JS exceptions
- title, URL, HTML size and a screenshot

## Usage

The container must be running first:

```bash
cd .. && make up
```

Run against rutube:

```bash
go run ./rutube -addr http://localhost:9222
```

Compare with a site the crawler handles fine:

```bash
go run ./rutube -addr http://localhost:9222 -url https://quotes.toscrape.com
```

## Flags

| Flag       | Default                 | Description                                   |
| ---------- | ----------------------- | --------------------------------------------- |
| `-addr`    | `http://localhost:9222` | CDP endpoint of the Tilion container          |
| `-url`     | `https://rutube.ru`     | Page to load                                  |
| `-timeout` | `30s`                   | Budget for load and stable waits              |
| `-output`  | `rutube.png`            | Screenshot output file                        |

Run `go run ./rutube -h` for the full list.

## Reading the output

A healthy page prints something like:

```text
t=0s        readyState=loading
t=1s        readyState=interactive
t=2s        readyState=complete
wait load:   ok in 1.9s
wait stable: ok in 5.2s
network: sent=40 responded=40 finished=40 failed=0 pending=0
```

If `wait stable` times out while `wait load` succeeded, look at
`still pending:` — every URL there started but never sent
`Network.loadingFinished`. A media stream (`*.mp4`, HLS segments) or a
long-polling XHR in that list is what blocks network-idle, and thus blocks
the crawler's "Wait page stable" forever.
