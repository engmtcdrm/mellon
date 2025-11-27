.PHONY: build runexe run test testv

build:
	@echo "Size before build:"; \
	ls -la |grep mellon; \
	ls -lh |grep mellon; \
	echo "\n\nSize after build:"; \
	CGO_ENABLED=0 go build --ldflags "-s -w"; \
	strip mellon; \
	ls -la |grep mellon; \
	ls -lh |grep mellon

runexe:
	@./mellon $(ARGS)

run:
	@go run . $(ARGS)

test:
	@go test -p=1 -timeout 30s ./...

testv:
	@go test -p=1 -timeout 30s -v ./...
