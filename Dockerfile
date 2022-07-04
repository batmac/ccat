FROM golang:1.18-alpine as build
WORKDIR /usr/src/app
COPY go.mod go.sum ./
RUN go mod download && go mod verify && apk add --no-cache build-base pkgconf curl-dev git
COPY . .
ENV CGO_ENABLED 1
RUN go version && ./build.sh

FROM alpine:20220328
RUN apk add --no-cache libcurl
COPY --from=build /usr/src/app/ccat /usr/bin/ccat

COPY "entrypoint.sh" "/entrypoint.sh"
CMD ["ccat"]
ENTRYPOINT ["/entrypoint.sh"]
HEALTHCHECK CMD /usr/bin/true
