# Valuation Rules Service

Bounded context: **Valuation Rules** — Decision rules binding service types to value and uncertainty envelopes.

Alignment: ISO 14008 monetary valuation posture. Decisions are `ALLOWED` / `DEGRADED` / `DENIED`, always
with the matched rule and a human-readable reason for audit. All data is synthetic.

| Item | Value |
|---|---|
| Port | `8083` (`HTTP_PORT`) |
| Base path | `/api/v1/valuation-rules` |
| Owns data | ValuationRule rules |
| Publishes | `valuation.activated` |

Run locally: `mvn -pl services/valuation-rules-service quarkus:dev`.
