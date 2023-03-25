FROM golang:alpine as builder

COPY . /go/src/github.com/Luzifer/mysqlapi
WORKDIR /go/src/github.com/Luzifer/mysqlapi

RUN set -ex \
 && apk add --update git \
 && go install \
      -ldflags "-X main.version=$(git describe --tags --always || echo dev)" \
      -mod=readonly \
      -modcacherw \
      -trimpath

FROM alpine:latest

LABEL maintainer "Knut Ahlers <knut@ahlers.me>"

RUN set -ex \
 && apk --no-cache add \
      ca-certificates

COPY --from=builder /go/bin/mysqlapi /usr/local/bin/mysqlapi

EXPOSE 3000

ENTRYPOINT ["/usr/local/bin/mysqlapi"]
CMD ["--"]

# vim: set ft=Dockerfile:
