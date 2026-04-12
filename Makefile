.PHONY: build test lint release release-dry licenses help

BINARY := nigopdf
VERSION ?=
TAG     := $(if $(filter v%,$(VERSION)),$(VERSION),v$(VERSION))

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-15s %s\n", $$1, $$2}'

build: ## Build the binary locally
	go build -o $(BINARY) .

test: ## Run all tests
	go test ./...

lint: ## Run go vet
	go vet ./...

licenses: ## Regenerate third-party license files
	./scripts/update-licenses.sh

release-dry: ## Dry-run GoReleaser locally (VERSION required)
	@if [ -z "$(VERSION)" ]; then echo "Usage: make release-dry VERSION=x.y.z"; exit 1; fi
	goreleaser release --snapshot --clean

release: ## Create tag and push to trigger release (VERSION required)
	@if [ -z "$(VERSION)" ]; then echo "Usage: make release VERSION=x.y.z"; exit 1; fi
	@if [ -n "$$(git status --porcelain)" ]; then echo "Error: working tree is not clean"; exit 1; fi
	git tag $(TAG)
	git push origin $(TAG)
	@echo "Tagged and pushed $(TAG). GitHub Actions will handle the release."
