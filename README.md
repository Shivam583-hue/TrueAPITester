# TrueAPITester

> Most API clients are Electron apps that eat 300MB of RAM, require an account, and sync your requests to someone else's cloud.
> This one is a single binary, opens in a millisecond, and keeps your data in a plain JSON file you own.

[![Go](https://img.shields.io/badge/Go-1.26-00ADD8?logo=go&logoColor=white)](https://go.dev)
[![Latest Release](https://img.shields.io/github/v/release/Shivam583-hue/TrueAPITester)](https://github.com/Shivam583-hue/TrueAPITester/releases/latest)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
[![Platform](https://img.shields.io/badge/platform-Linux%20%7C%20macOS%20%7C%20Windows-lightgrey)](#installation)

A fast, keyboard-driven API client that lives entirely in your terminal.

---

## Demo



https://github.com/user-attachments/assets/5a71ad80-b524-4171-913b-640a9555d375


---

## Why not Postman, Insomnia, or Bruno?

| | TrueAPITester | Postman | Insomnia | Bruno |
|---|---|---|---|---|
| **Runtime** | Single native binary | Electron (~300MB) | Electron (~300MB) | Electron (~200MB) |
| **Account required** | No | Yes (for sync) | Yes (for sync) | No |
| **Startup time** | Instant | 3-8 seconds | 3-8 seconds | 2-5 seconds |
| **Cloud sync of requests** | Never - local only | Optional, on by default | Optional, on by default | Never |
| **Stays in your terminal** | Yes | No | No | No |
| **Keyboard-driven** | Fully | Partial | Partial | Partial |

If you already live in the terminal - you `ssh` into servers, you run `git` from the command line, you have `nvim` or `tmux` open - switching to a GUI app just to fire off an API request breaks your flow.
TrueAPITester lets you stay where you are.

---

## Features

- **Collections** - organize requests in a sidebar, saved automatically to disk between sessions.
- **Full request editor** - method, URL, body, headers, query parameters, and auth (Bearer / Basic / API Key), each on its own tab.
- **Real HTTP requests** - sent with Go's standard `net/http`, no proxy or subprocess involved.
- **Readable responses** - syntax-highlighted JSON, raw body, response headers, and cookies, each on its own tab.
- **Run history** - every send is kept per request, so you can page back through past runs and compare status, timing, and size.
- **Context-aware help bar** - shows the keys relevant to whatever pane you're in, with a one-time full reference on first launch.

---

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

Or run directly without a separate build step:

```sh
go run ./cmd/APITester
```

---

## Usage

The window is split into a request sidebar, method/URL bar, request editor, and response pane.
`←`/`→` move focus between them; a short, context-aware hint bar at the bottom always shows what's available.
Press `?` at any time for the full keybinding reference.

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

| Key             | Action                                     |
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

| Key             | Action                         |
| --------------- | ------------------------------ |
| `tab`           | Switch result tab              |
| `↑` / `↓`       | Scroll                         |
| `pgup` / `pgdn` | Page-scroll                    |
| `[` / `]`       | Go to the previous / next run  |

### Global

| Key       | Action                                     |
| --------- | ------------------------------------------ |
| `ctrl+s`  | Send the active request                    |
| `ctrl+w`  | Save the collection to disk                |
| `?`       | Toggle the full help reference             |
| `ctrl+c`  | Quit (also saves the collection)           |

---

## Data storage

Collections are saved as human-readable JSON under your OS config directory:

- **Linux:** `~/.config/trueapitester/collection.json`
- **macOS:** `~/Library/Application Support/trueapitester/collection.json`
- **Windows:** `%AppData%\trueapitester\collection.json`

The file is loaded automatically on startup and saved automatically on quit.
You can also save manually at any time with `ctrl+w`.

The format is plain JSON - you can read it, back it up, diff it in git, or copy it between machines without any special tooling.
Your requests are never sent to any server.

If the file is corrupted (e.g. from a crash mid-write), TrueAPITester automatically renames it to `collection.json.bak` and starts fresh rather than crashing, so you never lose your previous state entirely.

Each request also stores up to 50 past responses as run history, so you can compare results across sessions without re-sending.

---

## Contributing

Contributions are welcome - see [CONTRIBUTING.md](CONTRIBUTING.md).

## License

[MIT](LICENSE)
