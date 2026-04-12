# GraphSentinel

**Structural code intelligence for detecting semantics-preserving transformations.**

GraphSentinel is a Go backend service for analyzing source code structure and detecting semantics-preserving transformations relevant to GNN robustness and code security research.

## Status

Minimal HTTP server is running: Chi router, graceful shutdown, and `GET /health`. Core analysis domain models, validation, and JSON shapes live in `pkg/models`. Submission and worker endpoints are next.

## Quickstart

```bash
make run
```

Optional environment:

| Variable | Default | Description |
|----------|---------|-------------|
| `HTTP_ADDR` | `:8080` | Listen address (`host:port` or `:port`) |
| `SHUTDOWN_TIMEOUT_SEC` | `15` | Graceful shutdown timeout in seconds |

Check health:

```bash
curl -s http://127.0.0.1:8080/health
```

Example response:

```json
{"status":"ok","service":"graphsentinel"}
```

## Layout

- `cmd/server` — application entrypoint
- `internal/` — API, config, ingestion, analyzers, detectors, workers, reports, store
- `pkg/models` — shared domain types
- `configs/`, `scripts/`, `deployments/`, `testdata/` — reserved for configuration and tooling

## License

Licensed under the [Apache License, Version 2.0](LICENSE).
