FROM golang:1.27-alpine as build
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

FROM alpine:3.24@sha256:28bd5fe8b56d1bd048e5babf5b10710ebe0bae67db86916198a6eec434943f8b
# hadolint ignore=DL3018
RUN apk upgrade --no-cache && apk add --no-cache libcurl tini
COPY "entrypoint.sh" "/entrypoint.sh"
COPY --from=build /usr/src/app/ccat /usr/bin/ccat
CMD ["ccat"]
ENTRYPOINT ["tini", "-wg", "--", "/entrypoint.sh"]
HEALTHCHECK CMD /usr/bin/true
