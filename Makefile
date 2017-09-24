
all: build

fmt:
	go fmt ./...

deps:
	glide install

install:
	env CGO_ENABLED=0 go install

clean:
	go clean -i

test:
	go test -v ./...

HAS_GLIDE := $(shell command -v glide;)
ifndef HAS_GLIDE
	go get -u github.com/Masterminds/glide
endif

.PHONY: fmt install clean test all
