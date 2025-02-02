BIN	= dedup
BINDIR	= $(HOME)/bin

$(BIN):	*.go
	CGO_ENABLED=0 go build

.PHONY: test
test:
	go vet
	staticcheck

.PHONY: clean
clean:
	go clean

.PHONY: gen
gen:
	@rm -f go.mod go.sum
	@go mod init $(BIN)
	@go mod tidy

.PHONY: install
install: $(BIN)
	@mkdir -p $(BINDIR)
	install -s -m 0755 $(BIN) $(BINDIR)
