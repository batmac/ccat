FROM golang:1.27-alpine@sha256:f92b6ef800e499660581efdabdf25d9d817a9d124eaf900924f0504e7e27e12d as build
WORKDIR /usr/src/app
COPY go.mod go.sum ./
# hadolint ignore=DL3018
RUN apk upgrade --no-cache \
    && apk add --no-cache build-base pkgconf curl-dev git bash
# populate the module cache in its own layer, invalidated only by go.mod/go.sum
RUN go mod download
COPY . .
ENV CGO_ENABLED=1
# zero-install mage: version pinned by go.mod instead of @latest
# hadolint ignore=DL3062
RUN go version && go run magefiles/mage.go buildFull

FROM alpine:3.24@sha256:be2656061b536c1883c4034790e21d0a09c5aaaf51968b89368a3c0371ea2c9a
# hadolint ignore=DL3018
RUN apk upgrade --no-cache && apk add --no-cache libcurl tini
COPY "entrypoint.sh" "/entrypoint.sh"
COPY --from=build /usr/src/app/ccat /usr/bin/ccat
CMD ["ccat"]
ENTRYPOINT ["tini", "-wg", "--", "/entrypoint.sh"]
HEALTHCHECK CMD /usr/bin/true
