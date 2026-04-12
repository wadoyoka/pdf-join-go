# Code Style & Conventions

## Go Conventions
- Standard Go formatting (gofmt)
- `go vet` for static analysis
- No global mutable state

## Project Patterns
- **cmd/ layer**: Thin wrappers — only flag parsing and user-facing output
- **internal/ layer**: All business logic, independently testable
- **New features**: Add cobra command in `cmd/`, implement logic in `internal/`
- **Error handling**: Return errors from `internal/`; `cmd/` handles cobra error flow
- **Error output**: Use `fmt.Fprintf(os.Stderr, ...)` for errors

## Testing
- Table-driven tests using `t.TempDir()` for isolation
- Tests create real PDFs via pdfcpu API (no mocks)
- Test files: `*_test.go` alongside source files

## Language
- Error messages: English
- README and user docs: Japanese

## Version
- Injected at build time via `-ldflags -X github.com/wadoyoka/nigopdf/cmd.version={{.Version}}`
