VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
LDFLAGS := -ldflags "-s -w -X agentmux/internal/commands.Version=$(VERSION)"
BIN     := agentmux

.PHONY: build install test lint fmt vet check clean

build:
	go build $(LDFLAGS) -o $(BIN) ./cmd/agentmux

install:
	go install $(LDFLAGS) ./cmd/agentmux

test:
	go test ./...

lint: fmt vet
	@echo "lint: ok"

fmt:
	@test -z "$$(gofmt -l .)" || (gofmt -d . && exit 1)

vet:
	go vet ./...

check: fmt vet test
	@echo "all checks passed"

clean:
	rm -f $(BIN)

smoke: build
	@echo "=== smoke test ==="
	./$(BIN) --version
	./$(BIN) list || true
	@echo "smoke: ok"
	@rm -f $(BIN)
