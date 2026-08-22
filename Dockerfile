FROM golang:1.27-alpine as build
WORKDIR /usr/src/app
COPY go.mod go.sum ./
# hadolint ignore=DL3018
RUN apk upgrade --no-cache \
    && apk add --no-cache build-base pkgconf curl-dev git bash \
    && go install github.com/magefile/mage@latest
# populate the module cache in its own layer, invalidated only by go.mod/go.sum
RUN go mod download
COPY . .
ENV CGO_ENABLED 1
RUN go version && mage buildFull

FROM alpine:20250108
# hadolint ignore=DL3018
RUN apk upgrade --no-cache && apk add --no-cache libcurl tini
COPY "entrypoint.sh" "/entrypoint.sh"
COPY --from=build /usr/src/app/ccat /usr/bin/ccat
CMD ["ccat"]
ENTRYPOINT ["tini", "-wg", "--", "/entrypoint.sh"]
HEALTHCHECK CMD /usr/bin/true
