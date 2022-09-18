FROM golang:1.18-alpine as deps

ADD go.mod /app/go.mod
WORKDIR /app

RUN go mod download




FROM deps as build
ARG TYPE

ADD . /app
WORKDIR /app

RUN time go build -o "/var/app" ./cmd/${TYPE} && \
    chmod a+x /var/app




FROM alpine

RUN addgroup -g 9999 -S user && \
    adduser -u 9999 -G user -S -H user

COPY --from=build /var/app /
ENTRYPOINT ["/app"]

USER user