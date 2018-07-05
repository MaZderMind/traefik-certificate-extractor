FROM golang:alpine as builder
RUN apk add --no-cache make git upx

WORKDIR /go/src/github.com/MaZderMind/traefik-certificate-extractor
COPY . .

RUN make binary \
    && upx --ultra-brute traefik-certificate-extractor

FROM scratch
COPY --from=builder /go/src/github.com/MaZderMind/traefik-certificate-extractor/traefik-certificate-extractor /

VOLUME /var/acmejson
ENTRYPOINT ["/traefik-certificate-extractor"]
CMD ["-acmejson=/var/acmejson/acme.json", "-target=/var/acmejson/certs/", "-watch"]
