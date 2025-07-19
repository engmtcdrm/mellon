.PHONY: build runexe run

build:
	echo "Size before build:"; ls -la |grep minno; ls -lh |grep minno; echo "\n\nSize after build:"; go build --ldflags "-s -w"; ls -la |grep minno; ls -lh |grep minno

runexe:
	./minno $(ARGS)

run:
	go run . $(ARGS)

test:
	go test ./...

testv:
	go test -v ./...
