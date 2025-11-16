GOBIN ?= $(HOME)/go/bin
BINARY := gojest
BUILDDIR := dist

.PHONY: install build clean

build:
	mkdir -p $(BUILDDIR)
	go build -o $(BUILDDIR)/$(BINARY) .

install: build
	mkdir -p $(GOBIN)
	rm -f $(GOBIN)/$(BINARY) || true
	cp $(BUILDDIR)/$(BINARY) $(GOBIN)/$(BINARY)
	@echo "Installed to $(GOBIN)/$(BINARY)"
	@make clean

clean:
	rm -rf $(BUILDDIR)


tag:
	@echo "Current version: $(VERSION)"
	@echo "Creating new version tag: v$(NEXT_VERSION)"
	git tag -a v$(NEXT_VERSION) -m "Release v$(NEXT_VERSION)"
	git push origin v$(NEXT_VERSION)
	@echo "✅ Tag v$(NEXT_VERSION) pushed to remote repository"

lint:	
	@command -v gopls >/dev/null 2>&1 || { \
		echo "🔧 Installing gopls..."; \
		go install golang.org/x/tools/gopls@latest; \
	}
	@echo "🔍 Running gopls check..."
	@gopls check $$(find . -name '*.go' -type f -not -path "./ent/*" -not -path "./vendor/*")
	@go test ./...