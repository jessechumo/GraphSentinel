# Stub: multi-stage build and non-root user will be finalized in a later commit.
FROM golang:1.23-alpine AS build
WORKDIR /src
COPY go.mod ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /graphsentinel ./cmd/server

FROM alpine:3.20
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=build /graphsentinel /app/graphsentinel
EXPOSE 8080
ENTRYPOINT ["/app/graphsentinel"]
