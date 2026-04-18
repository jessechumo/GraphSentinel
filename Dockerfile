FROM golang:1.23-alpine AS build
WORKDIR /src
COPY go.mod go.sum* ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -trimpath -o /graphsentinel ./cmd/server

FROM alpine:3.20
RUN apk add --no-cache ca-certificates wget \
	&& adduser -D -H -u 10001 graphsentinel
WORKDIR /app
COPY --from=build /graphsentinel /app/graphsentinel
USER 10001:10001
EXPOSE 8080
ENV HTTP_ADDR=:8080
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
	CMD wget -qO- http://127.0.0.1:8080/health >/dev/null || exit 1
ENTRYPOINT ["/app/graphsentinel"]
