# Provider lifecycle metrics

Provider lifecycle collectors track counters per provider without device IDs or
other high-cardinality values. The `TDLib` and `Whatsmeow` providers expose
`LifecycleMetrics()` for the application metrics adapter.

Each snapshot contains:

- `StartupAttempts`: actual provider connection attempts, including retries.
- `StartupSuccesses`: attempts that became usable connections.
- `StartupFailures`: failed attempts.
- `ActiveConnections`: saturating gauge of usable connections.
- `Disconnects`: active connections that completed cleanup.
- `LastFailureAt`: UTC timestamp of the most recent failed attempt.

Cleanup is idempotent at the metrics boundary: disconnecting a connection that
never became active does not make the active gauge negative or create a false
disconnect count. Device IDs are intentionally not retained in the collector.
