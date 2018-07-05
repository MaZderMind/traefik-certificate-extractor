.PHONY: default dependencies binary run container

GOFILES=traefik-certificate-extractor.go model.go

default: binary
binary: traefik-certificate-extractor

dependencies:
	go get

run: dependencies
	go run ${GOFILES}

traefik-certificate-extractor: ${GOFILES} dependencies
	CGO_ENABLED=0 GOOS=linux go build -a -ldflags="-s -w" -installsuffix cgo -o traefik-certificate-extractor .

container:
	docker build -t mazdermind/traefik-certificate-extractor:latest .
