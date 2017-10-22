
all: build

tools:
	which gometalinter || ( go get -u github.com/alecthomas/gometalinter && gometalinter --install )
	which glide || go get -u github.com/Masterminds/glide

lint:
	gometalinter --concurrency=1 --deadline=300s --vendor --disable-all \
		--enable=golint \
		--enable=vet \
		--enable=vetshadow \
		--enable=errcheck \
		--enable=structcheck \
		--enable=aligncheck \
		--enable=deadcode \
		--enable=ineffassign \
		--enable=dupl \
		--enable=gotype \
		--enable=varcheck \
		--enable=interfacer \
		--enable=goconst \
		--enable=megacheck \
		--enable=unparam \
		--enable=misspell \
		--enable=gas \
		--enable=goimports \
		./...

fmt:
	go fmt ./...

deps:
	glide install

build:
	env CGO_ENABLED=0 go build

install:
	env CGO_ENABLED=0 go install

clean:
	go clean -i

test:
	go test -v ./...

.PHONY: tools lint fmt install clean test all
