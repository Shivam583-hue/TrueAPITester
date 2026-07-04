# TrueAPITester

A fast, keyboard-driven API client that lives entirely in your terminal.

<!--
  Demo goes here. Once you have a recording, replace this comment with:
  ![demo](docs/demo.gif)
-->
<p align="center">
  <em>Demo coming soon</em>
</p>

## Features

- **Collections** — organize requests in a sidebar, saved automatically to disk between sessions.
- **Full request editor** — method, URL, body, headers, query parameters, and auth (Bearer / Basic / API Key), each on its own tab.
- **Real HTTP requests** — sent with Go's standard `net/http`, no proxy or subprocess involved.
- **Readable responses** — syntax-highlighted JSON, raw body, response headers, and cookies, each on its own tab.
- **Run history** — every send is kept per request, so you can page back through past runs and compare status, timing, and size.
- **Context-aware help bar** — shows the keys relevant to whatever pane you're in, with a one-time full reference on first launch.

## Installation

### Linux (x64)

```sh
curl -L https://github.com/Shivam583-hue/TrueAPITester/releases/latest/download/trueapitester-linux-amd64 -o apitester
chmod +x apitester
sudo mv apitester /usr/local/bin/
apitester
```

### Linux (ARM64)

```sh
curl -L https://github.com/Shivam583-hue/TrueAPITester/releases/latest/download/trueapitester-linux-arm64 -o apitester
chmod +x apitester
sudo mv apitester /usr/local/bin/
apitester
```

### macOS (Apple Silicon)

```sh
curl -L https://github.com/Shivam583-hue/TrueAPITester/releases/latest/download/trueapitester-macos-arm64 -o apitester
chmod +x apitester
sudo mv apitester /usr/local/bin/
apitester
```

### macOS (Intel)

```sh
curl -L https://github.com/Shivam583-hue/TrueAPITester/releases/latest/download/trueapitester-macos-amd64 -o apitester
chmod +x apitester
sudo mv apitester /usr/local/bin/
apitester
```

### Windows (PowerShell)

```powershell
curl.exe -L https://github.com/Shivam583-hue/TrueAPITester/releases/latest/download/trueapitester-windows-amd64.exe -o apitester.exe
move apitester.exe C:\Windows\System32\apitester.exe
apitester
```

### From source

Requires Go 1.26 or newer.

```sh
git clone https://github.com/Shivam583-hue/TrueAPITester.git
cd TrueAPITester
go build -o apitester ./cmd/APITester
./apitester
```

Or run it directly without a separate build step:

```sh
go run ./cmd/APITester
```

## Usage

The window is split into a request sidebar, method/URL bar, request editor, and response pane. `←`/`→` move focus between them; a short, context-aware hint bar at the bottom always shows what's available. Press `?` at any time for the full keybinding reference.

### Sidebar

| Key       | Action                  |
| --------- | ----------------------- |
| `n`       | New request             |
| `d`       | Delete selected request |
| `↑` / `↓` | Navigate requests       |

### Method / URL

| Key     | Action                    |
| ------- | ------------------------- |
| `m`     | Cycle HTTP method         |
| _type_  | Edit the URL              |
| `enter` | Move to the next pane     |

### Editor (Body / Headers / Query / Auth tabs)

| Key             | Action                                    |
| --------------- | ------------------------------------------ |
| `tab`           | Switch editor tab                          |
| _type_          | Edit body text (Body tab)                  |
| `n`             | Add a row (Headers / Query tabs)           |
| `d`             | Delete the selected row                    |
| `enter`         | Edit the selected row or field             |
| `esc`           | Cancel editing                             |
| `t`             | Cycle auth type (Auth tab)                 |
| `↑` / `↓`       | Move between rows, or scroll the body      |
| `pgup` / `pgdn` | Page-scroll the body                       |

### Result (Pretty / Raw / Headers / Cookies tabs)

| Key             | Action                        |
| --------------- | ----------------------------- |
| `tab`           | Switch result tab              |
| `↑` / `↓`       | Scroll                          |
| `pgup` / `pgdn` | Page-scroll                     |
| `[` / `]`       | Go to the previous / next run  |

### Global

| Key       | Action                                    |
| --------- | ------------------------------------------ |
| `ctrl+s`  | Send the active request                    |
| `ctrl+w`  | Save the collection to disk                |
| `?`       | Toggle the full help reference             |
| `ctrl+c`  | Quit (also saves the collection)           |

## Data storage

Collections (requests, editor state, and run history) are saved as JSON under your OS config directory, e.g. `~/.config/trueapitester/collection.json` on Linux. The file is loaded automatically on startup and saved automatically on exit.

## Contributing

Contributions are welcome — see [CONTRIBUTING.md](CONTRIBUTING.md).

## License

[MIT](LICENSE)
