BINARY := claude-menu
LDFLAGS := -s -w
PLATFORMS := darwin/amd64 darwin/arm64 windows/amd64 windows/arm64
DIST := dist
PREFIX ?= /usr/local

.PHONY: build install uninstall release clean

## build: compile a binary for the current OS/arch
build:
	go build -ldflags "$(LDFLAGS)" -o $(BINARY) .

## install: build and copy to $(PREFIX)/bin (macOS/Linux; use install.ps1 on Windows)
install: build
	install -d $(PREFIX)/bin
	install -m 0755 $(BINARY) $(PREFIX)/bin/$(BINARY)
	@echo "Installed $(PREFIX)/bin/$(BINARY)"

## uninstall: remove the installed binary
uninstall:
	rm -f $(PREFIX)/bin/$(BINARY)

## release: cross-compile release binaries for all platforms into $(DIST)/
release: clean
	@mkdir -p $(DIST)
	@for p in $(PLATFORMS); do \
	  os=$${p%/*}; arch=$${p#*/}; \
	  ext=""; [ "$$os" = "windows" ] && ext=".exe"; \
	  out=$(DIST)/$(BINARY)_$${os}_$${arch}$$ext; \
	  echo "building $$out"; \
	  GOOS=$$os GOARCH=$$arch CGO_ENABLED=0 go build -ldflags "$(LDFLAGS)" -o $$out . || exit 1; \
	done
	@cd $(DIST) && shasum -a 256 * > checksums.txt
	@echo "Artifacts in $(DIST)/"

## clean: remove build outputs
clean:
	rm -rf $(DIST) $(BINARY)
