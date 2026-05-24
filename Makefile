VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
LDFLAGS := -ldflags "-s -w -X agentmux/internal/commands.Version=$(VERSION)"
BIN     := agentmux
PREFIX  ?= /usr/local

.PHONY: build install uninstall test lint fmt vet check clean smoke

build:
	go build $(LDFLAGS) -o $(BIN) ./cmd/agentmux

install: build
	install -d $(PREFIX)/bin
	install -m 755 $(BIN) $(PREFIX)/bin/$(BIN)
	ln -sf $(BIN) $(PREFIX)/bin/atmux
	@echo "Installed $(PREFIX)/bin/agentmux"
	@echo "Installed $(PREFIX)/bin/atmux -> agentmux"

uninstall:
	rm -f $(PREFIX)/bin/$(BIN) $(PREFIX)/bin/atmux

test:
	go test -race ./...

lint: fmt vet
	@echo "lint: ok"

fmt:
	@test -z "$$(gofmt -l .)" || (gofmt -d . && exit 1)

vet:
	go vet ./...

check: fmt vet test build
	@echo "all checks passed"

clean:
	rm -f $(BIN)

smoke: build
	@echo "=== smoke test ==="
	@./$(BIN) --version
	@./$(BIN) list || true
	@echo "=== smoke: ok ==="
	@rm -f $(BIN)
