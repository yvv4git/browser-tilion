# pixelscan_botcheck

Opens the PixelScan bot-check page through a local HTTP proxy, waits for a
configurable time, saves a timestamped screenshot into `examples/.screenshots/`,
then closes the page and the browser.

Launches its own browser via rod launcher (not the Tilion container) because
the proxy must be applied at browser launch time.

## Usage

Proxy must be running first (e.g. `http://localhost:8138`).

Default (10s wait, screenshot to `../.screenshots/`):

```bash
go run ./pixelscan_botcheck -proxy http://localhost:8138
```

Custom wait:

```bash
go run ./pixelscan_botcheck -proxy http://localhost:8138 -wait 15s
```

Custom output path:

```bash
go run ./pixelscan_botcheck -proxy http://localhost:8138 -output ../.screenshots/custom.png
```

## Flags

| Flag       | Default                                       | Description                          |
| ---------- | --------------------------------------------- | ------------------------------------ |
| `-proxy`   | `http://localhost:8138`                       | HTTP proxy URL                       |
| `-url`     | `https://pixelscan.net/bot-check`             | Page to load                         |
| `-wait`    | `10s`                                         | Time to wait before screenshot       |
| `-output`  | `../.screenshots/pixelscan_botcheck_<ts>.png` | Output file (timestamped by default) |

Run `go run ./pixelscan_botcheck -h` for the full list.
