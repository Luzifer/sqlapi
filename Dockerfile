FROM golang:1.25-alpine@sha256:b6ed3fd0452c0e9bcdef5597f29cc1418f61672e9d3a2f55bf02e7222c014abd AS builder

COPY . /go/src/github.com/Luzifer/sqlapi
WORKDIR /go/src/github.com/Luzifer/sqlapi

RUN set -ex \
 && apk add --update git \
 && go install \
      -ldflags "-X main.version=$(git describe --tags --always || echo dev)" \
      -mod=readonly \
      -modcacherw \
      -trimpath


FROM alpine:3.22@sha256:4bcff63911fcb4448bd4fdacec207030997caf25e9bea4045fa6c8c44de311d1

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
