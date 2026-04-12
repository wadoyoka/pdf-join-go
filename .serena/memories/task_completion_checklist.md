# Task Completion Checklist

When a task is completed, run the following:

1. **Format**: `gofmt -w .` (or ensure code is formatted)
2. **Vet**: `go vet ./...`
3. **Test**: `go test ./...`
4. **Verify build**: `go build -o nigopdf .`

If `go.mod` or `go.sum` were modified, lefthook pre-commit will automatically regenerate:
- `THIRD_PARTY_LICENSES/`
- `CREDITS.txt`
- `cmd/credits.txt`
