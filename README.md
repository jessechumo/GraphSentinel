# GraphSentinel

**Structural code intelligence for detecting semantics-preserving transformations.**

GraphSentinel is a Go backend service for analyzing source code structure and detecting semantics-preserving transformations relevant to GNN robustness and code security research.

## Status

HTTP server with `GET /health`, **`POST /analyze`**, and **`GET /analysis/{id}`**. Submissions are queued in `internal/store`, processed by a configurable worker pool (`internal/workers`), and analyzed with MVP identifier-renaming, dead-code, and control-flow detectors. Errors and access logs use a consistent JSON shape on stdout for production-style observability.

## Quickstart

```bash
make run
```

Validate before committing:

```bash
make check
```

Optional environment:

| Variable | Default | Description |
|----------|---------|-------------|
| `HTTP_ADDR` | `:8080` | Listen address (`host:port` or `:port`) |
| `LOG_LEVEL` | `info` | Log level: `debug`, `info`, `warn`, `error` |
| `SHUTDOWN_TIMEOUT_SEC` | `15` | Graceful shutdown timeout in seconds |
| `WORKER_COUNT` | `2` | Concurrent analysis workers |
| `WORKER_QUEUE_SIZE` | `256` | Buffered job queue depth for the worker pool |
| `READ_TIMEOUT_SEC` | `30` | HTTP server read timeout |
| `WRITE_TIMEOUT_SEC` | `60` | HTTP server write timeout |
| `IDLE_TIMEOUT_SEC` | `120` | HTTP server idle connection timeout |

Structured logs are JSON lines written to stdout (`log/slog`). Example fields: `http_request` (per request) and `api_error` (4xx/5xx responses).

See `configs/.env.example` for a copy-paste template.

Check health:

```bash
curl -s http://127.0.0.1:8080/health
```

Example response:

```json
{"status":"ok","service":"graphsentinel"}
```

Submit code for analysis (returns `202 Accepted` with a queued job id):

```bash
curl -s -X POST http://127.0.0.1:8080/analyze \
  -H 'Content-Type: application/json' \
  -d '{"language":"c","code":"int main(){return 0;}"}'
```

Example response:

```json
{"status":"queued","analysis_id":"<hex-id>"}
```

Fetch status or the completed report (poll until `status` is `completed` or `failed`):

```bash
curl -s http://127.0.0.1:8080/analysis/<hex-id>
```

Validation error (standard error shape):

```bash
curl -s -X POST http://127.0.0.1:8080/analyze \
  -H 'Content-Type: application/json' \
  -d '{"language":"c","code":""}'
```

Example response:

```json
{"error":"code is required","code":"VALIDATION_ERROR"}
```

Run in Docker:

```bash
docker build -t graphsentinel:dev .
docker run --rm -p 8080:8080 graphsentinel:dev
```

## Layout

- `cmd/server` — application entrypoint
- `internal/` — API, config, ingestion, analyzers, detectors, workers, reports, store
- `pkg/models` — shared domain types
- `configs/`, `scripts/`, `deployments/`, `testdata/` — reserved for configuration and tooling

## License

Licensed under the [Apache License, Version 2.0](LICENSE).
