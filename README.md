# GraphSentinel

**Structural code intelligence for detecting semantics-preserving transformations.**

GraphSentinel is a Go backend service for analyzing source code structure and detecting semantics-preserving transformations relevant to GNN robustness and code security research.

## Status

HTTP server with `GET /health`, **`POST /analyze`**, and **`GET /analysis/{id}`**. Submissions are queued in `internal/store`, processed by a configurable worker pool (`internal/workers`), and completed with a baseline stub report (`internal/reports`) until the detector commits land.

## Quickstart

```bash
make run
```

Optional environment:

| Variable | Default | Description |
|----------|---------|-------------|
| `HTTP_ADDR` | `:8080` | Listen address (`host:port` or `:port`) |
| `SHUTDOWN_TIMEOUT_SEC` | `15` | Graceful shutdown timeout in seconds |
| `WORKER_COUNT` | `2` | Concurrent analysis workers |

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

## Layout

- `cmd/server` — application entrypoint
- `internal/` — API, config, ingestion, analyzers, detectors, workers, reports, store
- `pkg/models` — shared domain types
- `configs/`, `scripts/`, `deployments/`, `testdata/` — reserved for configuration and tooling

## License

Licensed under the [Apache License, Version 2.0](LICENSE).
