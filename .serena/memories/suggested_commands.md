# Suggested Commands

## Build & Run
```bash
# Build the binary
go build -o nigopdf .

# Run directly
go run .

# Example usage
./nigopdf merge file1.pdf file2.pdf --output merged.pdf
./nigopdf split input.pdf --parts 3 --output ./output/
./nigopdf split input.pdf --max-size 10MB --output ./output/
```

## Testing
```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run specific package tests
go test ./internal/merger/
go test ./internal/splitter/
```

## Code Quality
```bash
# Lint / vet
go vet ./...

# Format
gofmt -w .
```

## Git & Release
```bash
# Git hooks are managed by lefthook
# Pre-commit: auto-regenerates THIRD_PARTY_LICENSES/, CREDITS.txt, cmd/credits.txt

# Release: tag and push
git tag v1.x.x
git push origin v1.x.x
# GitHub Actions + GoReleaser handles the rest
```

## System Utils (Darwin)
```bash
git, ls, cd, grep, find, cat, head, tail
# Note: macOS uses BSD versions of these tools
```
