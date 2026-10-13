# screenshot

Opens a page in the Tilion container and saves screenshots:

- viewport screenshot (default)
- full-page screenshot with `-full`

## Usage

The container must be running first:

```bash
cd .. && make up
```

Viewport screenshot:

```bash
go run ./screenshot -addr http://localhost:9222 -url https://www.wikipedia.org
```

Full-page screenshot:

```bash
go run ./screenshot -addr http://localhost:9222 -url https://www.wikipedia.org -full -output wikipedia_full.png
```

## Flags

| Flag      | Default                     | Description                          |
| --------- | --------------------------- | ------------------------------------ |
| `-addr`   | `http://localhost:9222`     | CDP endpoint of the Tilion container |
| `-url`    | `https://www.wikipedia.org` | Page to load                         |
| `-output` | `screenshot.png`            | Output file                          |
| `-full`   | `false`                     | Capture full page instead of viewport |

Run `go run ./screenshot -h` for the full list.
