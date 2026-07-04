# Contributing to TrueAPITester

Thanks for taking a look at this project. It's a small, focused tool, so the goal for contributions is to keep it that way — simple, readable, and easy to reason about.

## Getting set up

You'll need Go 1.26+.

```sh
git clone https://github.com/Shivam583-hue/TrueAPITester.git
cd TrueAPITester
go run ./cmd/APITester
```

Before sending a change, make sure it passes the usual checks:

```sh
gofmt -l .        # should print nothing
go vet ./...
go build ./...
```

There's no test suite yet; if you add one for a package, make sure `go test ./...` passes too.

## Project layout

| Path                    | What lives there                                                       |
| ----------------------- | ------------------------------------------------------------------------ |
| `cmd/APITester`         | Entry point — wires up the Bubble Tea program                            |
| `internal/model`        | The TUI itself: layout, key handling, rendering (Bubble Tea Model/Update/View) |
| `internal/store`        | The request collection: types, CRUD, run history, JSON persistence — no UI dependency |
| `internal/httpClient`   | Sends the actual HTTP requests                                           |
| `internal/styles`       | Shared Lipgloss styles and layout helpers                                |

`internal/store` is intentionally decoupled from `internal/model` — it doesn't import Bubble Tea and can be tested or reused on its own.

## Making changes

- Keep pull requests focused on one thing. Unrelated cleanup makes a change harder to review and easier to revert by accident.
- Match the existing style: no unnecessary comments, no abstractions for cases that don't exist yet, `gofmt`-formatted.
- If you're adding a keybinding, update the relevant `key.Binding` in `internal/model/keymap.go` so it shows up in the help bar, and add it to the README's keybinding tables.
- If you're changing anything persisted to `collection.json` (in `internal/store`), think about what happens when an older collection file is loaded — old fields should degrade gracefully rather than crash.

## Reporting bugs / requesting features

Open a GitHub issue with:

- What you expected to happen vs. what actually happened
- Steps to reproduce (a specific request/response shape helps a lot for HTTP-related bugs)
- Your OS and terminal emulator, if the issue looks visual (rendering, colors, layout)

## Submitting a pull request

1. Fork the repo and create a branch off `main`.
2. Make your change, verify it with the checks above.
3. Open a PR describing what changed and why. Screenshots or a terminal recording are especially helpful for anything touching the UI.
