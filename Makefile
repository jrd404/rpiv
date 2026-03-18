VERSION ?= dev
BINARY  := rpiv
LDFLAGS := -buildvcs=false -ldflags "-X main.version=$(VERSION)"

.PHONY: build install test lint fmt vet check clean sync-assets

## Build

build: sync-assets
	go build $(LDFLAGS) -o bin/$(BINARY) ./cmd/rpiv

install: sync-assets
	go install $(LDFLAGS) ./cmd/rpiv/

## Sync plugin/skill assets into internal/assets/ for go:embed
sync-assets:
	@rsync -a --delete plugins/ internal/assets/plugins/
	@rsync -a --delete skills/ internal/assets/skills/

## Test & Lint

test:
	go test ./...

fmt:
	gofmt -w .

vet:
	go vet ./...

lint:
	golangci-lint run

check: fmt vet lint test

## Completions

completion-bash:
	@go run $(LDFLAGS) ./cmd/rpiv/ completion bash

completion-zsh:
	@go run $(LDFLAGS) ./cmd/rpiv/ completion zsh

install-completions: build
	@echo "Installing bash completion..."
	@mkdir -p /etc/bash_completion.d 2>/dev/null || true
	@cp completions/rpiv.bash /etc/bash_completion.d/rpiv 2>/dev/null || echo "  (skipped /etc/bash_completion.d - use 'complete -C rpiv\\ __completer rpiv' instead)"
	@echo "Installing zsh completion..."
	@mkdir -p ~/.zsh/completions 2>/dev/null || true
	@cp completions/rpiv.zsh ~/.zsh/completions/_rpiv

## Clean

clean:
	rm -rf bin/
	rm -rf internal/assets/plugins/ internal/assets/skills/
