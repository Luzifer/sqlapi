FROM golang:1.27.1-alpine@sha256:cf6fca6641884b8433441b2b0652976f975e1d0fdd26d177eaaf8596087f3125 AS builder

COPY . /go/src/github.com/Luzifer/sqlapi
WORKDIR /go/src/github.com/Luzifer/sqlapi

RUN set -ex \
 && apk add --update git \
 && go install \
      -ldflags "-X main.version=$(git describe --tags --always || echo dev)" \
      -mod=readonly \
      -modcacherw \
      -trimpath


FROM alpine:3.24.1@sha256:28bd5fe8b56d1bd048e5babf5b10710ebe0bae67db86916198a6eec434943f8b

LABEL org.opencontainers.image.authors="Knut Ahlers <knut@ahlers.me>" \
      org.opencontainers.image.source="https://git.luzifer.io/luzifer/sqlapi"

RUN set -ex \
 && apk --no-cache add \
      ca-certificates \
      tzdata

COPY --from=builder /go/bin/sqlapi /usr/local/bin/sqlapi

EXPOSE 3000

ENTRYPOINT ["/usr/local/bin/sqlapi"]
CMD ["--"]

# vim: set ft=Dockerfile:
