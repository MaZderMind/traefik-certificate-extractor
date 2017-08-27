FROM scratch
ADD traefik-certificate-extractor /

VOLUME /var/acmejson
CMD ["/traefik-certificate-extractor", "-acmejson=/var/acmejson/acme.json", "-target=/var/acmejson/certs/", "-watch"]
