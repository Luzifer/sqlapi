FROM golang:1.27.4-alpine@sha256:0a5f7fe796599ff487605987432e7619ae45beeb2fba4c2ce0b907879932aebf AS builder

COPY . /go/src/github.com/Luzifer/sqlapi
WORKDIR /go/src/github.com/Luzifer/sqlapi

RUN set -ex \
 && apk add --update git \
 && go install \
      -ldflags "-X main.version=$(git describe --tags --always || echo dev)" \
      -mod=readonly \
      -modcacherw \
      -trimpath


FROM alpine:3.24.4@sha256:f59abe656b692d901e184be3445e084786602108599d2a51d229c0f5a48c71bf

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
