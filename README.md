# GraphSentinel

**Structural code intelligence for detecting semantics-preserving transformations.**

GraphSentinel is a Go backend service for analyzing source code structure and detecting semantics-preserving transformations relevant to GNN robustness and code security research.

## Status

Repository skeleton initialized. HTTP API, async analysis, and detectors land in upcoming commits.

## Quickstart (skeleton)

```bash
make run
```

## Layout

- `cmd/server` — application entrypoint
- `internal/` — API, config, ingestion, analyzers, detectors, workers, reports, store
- `pkg/models` — shared domain types
- `configs/`, `scripts/`, `deployments/`, `testdata/` — reserved for configuration and tooling

## License

Specify your license here.
