# CLAUDE.md

## Task Routing

Always use the `/tasks-router` skill for every user instruction. This skill analyzes the intent, selects appropriate skills and subagents, creates an execution plan, and executes after user approval. Do not skip this step even for seemingly simple tasks.

## Project Overview

nigopdf is a Go CLI tool that simplifies PDF merging and splitting. It wraps pdfcpu with a user-friendly cobra-based CLI.

## Architecture

```
cmd/           # Thin CLI wrappers (cobra commands)
  root.go      # Root command, --version, --credits
  merge.go     # merge subcommand
  split.go     # split subcommand
internal/      # Core business logic
  merger/      # PDF collection and merging
  splitter/    # PDF splitting (by parts, by max size)
```

- `cmd/` handles flag parsing and user-facing output only
- `internal/` contains all logic and is independently testable
- No shared state between packages

## Build & Test

```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run specific package tests
go test ./internal/merger/
go test ./internal/splitter/

# Build
go build -o nigopdf .

# Lint
go vet ./...
```

## Dependencies

- `github.com/pdfcpu/pdfcpu` — PDF manipulation engine
- `github.com/spf13/cobra` — CLI framework
- Go 1.25.5

## Conventions

- Tests are table-driven using `t.TempDir()` for isolation
- Test files create real PDFs via pdfcpu API (no mocks)
- Error messages are in English; README and user docs are in Japanese
- Version is injected at build time via `-ldflags -X`
- All new CLI features follow the pattern: add cobra command in `cmd/`, implement logic in `internal/`

## Git Hooks (lefthook)

- Pre-commit: auto-regenerates license files when `go.mod`/`go.sum` change
- Generated files: `THIRD_PARTY_LICENSES/`, `CREDITS.txt`, `cmd/credits.txt`

## Release Process

- Tag with `v*` and push to trigger GitHub Actions
- GoReleaser builds darwin/linux (amd64/arm64) binaries
- Publishes to Homebrew tap `wadoyoka/homebrew-tap`
- Requires `GITHUB_TOKEN` and `HOMEBREW_TAP_GITHUB_TOKEN`

## Serena MCP

コードの読み書きには Serena MCP ツールを優先的に使用すること。

- シンボルの概要把握には `get_symbols_overview` を使う
- シンボルの検索には `find_symbol` を使う
- シンボルの編集には `replace_symbol_body`、`insert_before_symbol`、`insert_after_symbol` を使う
- リファクタリングには `rename_symbol`、`safe_delete_symbol` を使う
- ファイル全体を読む前に、まずシンボル単位で必要な情報だけを取得する

## Tool Usage

- ファイル内容の検索には Bash の `grep` / `rg` ではなく、専用の `Grep` ツールを使うこと
- Go module cache (`~/go/pkg/mod/`) など外部パスの検索も `Grep` ツールで行う

## Code Style

- Keep `cmd/` functions thin — delegate to `internal/`
- Use `fmt.Fprintf(os.Stderr, ...)` for error output
- Return errors from `internal/` functions; let `cmd/` handle `cobra.Command` error flow
- No global mutable state
