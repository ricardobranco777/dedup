BIN	= dedup
BINDIR	= $(HOME)/bin

GO	?= go
DOCKER	?= podman

# https://github.com/golang/go/issues/64875
arch := $(shell uname -m)
ifeq ($(arch),s390x)
CGO_ENABLED := 1
else
CGO_ENABLED ?= 0
endif

$(BIN): *.go
	CGO_ENABLED=$(CGO_ENABLED) $(GO) build -trimpath -ldflags="-s -w -buildid=" -buildmode=pie

.PHONY: build
build:
	image=$$( $(DOCKER) build -q . ) && \
	container=$$( $(DOCKER) create $$image ) && \
	$(DOCKER) cp $$container:/usr/local/bin/$(BIN) . && \
	$(DOCKER) rm -vf $$container && \
	$(DOCKER) rmi $$image

.PHONY: test
test:
	$(GO) vet
	staticcheck
	gofmt -s -l .

.PHONY: lint
lint:
	golangci-lint run ./...

.PHONY: clean
clean:
	$(GO) clean

.PHONY: gen
gen:
	$(RM) go.mod go.sum
	$(GO) mod init github.com/ricardobranco777/$(BIN)
	$(GO) mod tidy

.PHONY: install
install: $(BIN)
	@mkdir -p $(BINDIR)
	install -s -m 0755 $(BIN) $(BINDIR)
