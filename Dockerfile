FROM golang:1.21-alpine as build
WORKDIR /usr/src/app
COPY go.mod go.sum ./
# hadolint ignore=DL3018
RUN apk upgrade --no-cache \
    && apk add --no-cache build-base pkgconf curl-dev git bash \
    && go install github.com/magefile/mage@latest
COPY . .
ENV CGO_ENABLED 1
RUN go version && mage buildFull

FROM alpine:20231219
# hadolint ignore=DL3018
RUN apk upgrade --no-cache && apk add --no-cache libcurl tini
COPY "entrypoint.sh" "/entrypoint.sh"
COPY --from=build /usr/src/app/ccat /usr/bin/ccat
CMD ["ccat"]
ENTRYPOINT ["tini", "-wg", "--", "/entrypoint.sh"]
HEALTHCHECK CMD /usr/bin/true
