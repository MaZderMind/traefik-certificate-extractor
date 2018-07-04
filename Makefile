.PHONY: default binary run container push

GOFILES=traefik-certificate-extractor.go model.go

default: binary
binary: traefik-certificate-extractor

run:
	go run ${GOFILES}

traefik-certificate-extractor: ${GOFILES}
	CGO_ENABLED=0 GOOS=linux go build -a -ldflags="-s -w" -installsuffix cgo -o traefik-certificate-extractor .

container: binary
	docker build -t mazdermind/traefik-certificate-extractor:latest .

push: container
	docker push mazdermind/traefik-certificate-extractor:latest
