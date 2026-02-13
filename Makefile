.PHONY: test test-run generate tidy build

# Run all tests
test:
	go test -v ./...

# Run a specific test by pattern: make test-run RUN=042
test-run:
	go test -v ./datetime -run '$(RUN)'

# Regenerate yacc parser from grammar
generate:
	cd datetime && go generate ./...

# Clean up go.mod
tidy:
	go mod tidy

# Build all packages
build:
	go build ./...
