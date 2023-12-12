
FROM golang:1.21.5-alpine3.19 AS builder
ENV APP_NAME ponder
ENV WORKDIR /app
WORKDIR $WORKDIR
COPY . .
RUN go mod download
RUN go build -o /$APP_NAME

## Deploy
FROM alpine:3.19.0
ENV APP_NAME ponder
WORKDIR /
COPY --from=builder /$APP_NAME /$APP_NAME
ENTRYPOINT ["/ponder"]
