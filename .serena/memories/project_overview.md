# nigopdf - Project Overview

## Purpose
nigopdf is a Go CLI tool for PDF merging and splitting. It wraps pdfcpu with a user-friendly cobra-based CLI.

## Tech Stack
- **Language**: Go 1.25.5
- **CLI framework**: github.com/spf13/cobra
- **PDF engine**: github.com/pdfcpu/pdfcpu v0.11.1
- **Build/Release**: GoReleaser v2, GitHub Actions
- **Git hooks**: lefthook (pre-commit: auto-regenerate license files)
- **Platforms**: darwin/linux (amd64/arm64)

## Architecture
```
main.go            # Entry point, calls cmd.Execute()
cmd/               # Thin CLI wrappers (cobra commands)
  root.go          # Root command, --version, --credits flags
  merge.go         # merge subcommand (flags: --output, --dry-run, --recursive)
  split.go         # split subcommand (flags: --parts, --max-size, --output)
  credits.txt      # Embedded credits text
internal/          # Core business logic (independently testable)
  merger/          # PDF collection (flat/recursive) and merging
    merger.go      # isPDF, CollectPDFs, collectPDFsFlat, collectPDFsRecursive, Merge
    merger_test.go
  splitter/        # PDF splitting (by parts, by max size)
    splitter.go    # ParseSize, outputFileName, SplitByParts, SplitByMaxSize, readContext
    splitter_test.go
scripts/           # Helper scripts (update-licenses.sh)
```

## Key Design Principles
- `cmd/` handles only flag parsing and user-facing output
- `internal/` contains all logic and is independently testable
- No shared state between packages
- Errors returned from `internal/`; `cmd/` handles cobra error flow
- `fmt.Fprintf(os.Stderr, ...)` for error output
- Version injected at build time via `-ldflags -X`
