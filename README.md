traefik-certificate-extractor
=============================

A small utility which monitors a traefik-managed acme.json and extracts the plain certificate-files from it.

This is a fork based on the [excelent work of MaZderMind](https://github.com/MaZderMind/traefik-certificate-extractor).

For each domain and it extracts
 - `fullchain` (cert + intermediate)
 - `privkey` (private key)
 - `all` (private key + cert + intermediate)

For SANs it creates symlinks to the main domain's files.

It can be used in one-shot mode or in watch-mode, monitoring changes to the acme.json.
The tool creates and writes files but never removes old files or directories.

```
Usage of ./traefik-certificate-extractor:
  -acmejson string
    	path of the acme.json-file
  -target string
    	directory where the certificates should be extracted to
  -watch
    	should the extractor-tool keep watching the acme.json-file and rewrite the certificates
```

Docker-Container
----------------
[sbruder/traefik-certificate-extractor](https://hub.docker.com/r/sbruder/traefik-certificate-extractor/)

Expects `acme.json` in `/var/acmejson/acme.json`, writes certs to `/var/acmejson/certs/`
