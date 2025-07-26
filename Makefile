.PHONY: build runexe run

build:
	echo "Size before build:"; ls -la |grep mellon; ls -lh |grep mellon; echo "\n\nSize after build:"; go build --ldflags "-s -w"; ls -la |grep mellon; ls -lh |grep mellon

runexe:
	./mellon $(ARGS)

run:
	go run . $(ARGS)

test:
	go test ./...

testv:
	go test -v ./...
