FROM golang:1.19-alpine as build
WORKDIR /usr/src/app
COPY go.mod go.sum ./
RUN go mod download && go mod verify && apk add --no-cache build-base pkgconf curl-dev git
COPY . .
ENV CGO_ENABLED 1
RUN go version && ./build.sh

FROM alpine:20220715
RUN apk add --no-cache libcurl
COPY "entrypoint.sh" "/entrypoint.sh"
COPY --from=build /usr/src/app/ccat /usr/bin/ccat
CMD ["ccat"]
ENTRYPOINT ["/entrypoint.sh"]
HEALTHCHECK CMD /usr/bin/true
