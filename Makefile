GO ?= go

.PHONY: build test check acceptance provider-build

build:
	$(GO) build ./...

test:
	$(GO) test ./...

provider-build:
	$(GO) build -o dist/terraform-provider-open-unifi ./cmd/tfprovider

check:
	@out=$$(for f in $$(git ls-files '*.go'); do [ -f "$$f" ] || continue; gofmt -l "$$f"; done | sort -u); if [ -n "$$out" ]; then \
		echo "gofmt needed on:"; echo "$$out"; exit 1; \
	fi
	$(GO) vet ./...
	$(GO) test -count=1 ./...

acceptance:
	TF_ACC=1 $(GO) test -tags acceptance -count=1 ./acceptance -v
