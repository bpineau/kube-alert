
all: build

fmt:
	go fmt ./...

deps:
	which glide || go get -u github.com/Masterminds/glide
	glide install

build:
	env CGO_ENABLED=0 go build

install:
	env CGO_ENABLED=0 go install

clean:
	go clean -i

test:
	go test -v ./...

.PHONY: fmt install clean test all
