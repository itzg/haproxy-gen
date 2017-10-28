
GLIDE := ${GOPATH}/bin/glide
GORELEASER := ${GOPATH}/bin/goreleaser

default: dependencies test install

dependencies: ${GLIDE}
	glide install

test:
		go test $(shell glide nv)

build:
	go build

install:
	go install

snapshot: ${GORELEASER} .goreleaser.yml
	${GORELEASER} --snapshot

release: ${GORELEASER} .goreleaser.yml
	${GORELEASER}

${GLIDE}:
	curl https://glide.sh/get | sh

${GORELEASER}:
	go get github.com/goreleaser/goreleaser

.PHONY : test dependencies build install release
