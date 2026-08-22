FROM golang:1.27-alpine@sha256:4c9fe60190a2a3350ddc51de80d0224b8a6698d12bdfc999fee45ea9d6c46dbc as build
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
