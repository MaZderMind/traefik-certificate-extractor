FROM golang:alpine as builder
RUN apk add --no-cache make git

WORKDIR /go/src/github.com/MaZderMind/traefik-certificate-extractor
COPY . .

RUN make binary

FROM scratch
COPY --from=builder /go/src/github.com/MaZderMind/traefik-certificate-extractor/traefik-certificate-extractor /

VOLUME /var/acmejson
CMD ["/traefik-certificate-extractor", "-acmejson=/var/acmejson/acme.json", "-target=/var/acmejson/certs/", "-watch"]
