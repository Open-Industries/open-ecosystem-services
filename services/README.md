# Services

One independently deployable Quarkus microservice per bounded context. Every service
follows the same hexagonal layout, exposes `/openapi`, `/q/health` and `/q/metrics`,
returns RFC 9457 problem details, and is configured only through environment variables
(Twelve-Factor). Lifecycle transitions take no body; `@Consumes` appears only on
endpoints that read one.

**Implementation order** (see `../docs/implementation-roadmap.md`):
Ecosystem Assets → Service Credits → Valuation Rules → Open APIs.

| Service | Port | Alignment |
|---|---|---|
| `ecosystem-asset-service` | 8081 | UN SEEA EA ecosystem extent accounts |
| `service-credit-service` | 8082 | UN SEEA EA ecosystem services accounts |
| `valuation-rules-service` | 8083 | ISO 14008 monetary valuation posture |
| `monitoring-station-service` | 8084 | UN SEEA EA condition accounts (observation sources) |
| `notification-service` | 8085 | Platform support (shared pattern) |

Run the pilot service:

```bash
mvn -pl services/ecosystem-asset-service quarkus:dev
# OpenAPI UI: http://localhost:8081/docs — health: http://localhost:8081/q/health
```
