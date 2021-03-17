all: build

build:
	go build

clean:
	rm -f imgchroot coverage.out go.sum
	rm -rf dist/ site/

docs:
	mkdocs build

lint:
	golangci-lint run

release: clean
	goreleaser

release-snapshot: clean
	goreleaser --rm-dist --skip-publish --snapshot

test:
	go test -cover -coverprofile=coverage.out  ./pkg/...

test-coverage: test
	go tool cover -html=coverage.out

.PHONY: all build clean docs lint release release-snapshot test test-coverage
