# Open Ecosystem Services — Reference Ecosystem Services Core Platform

A **standard, generic ecosystem services and natural capital core platform** that any organization in this vertical
can deploy and adapt: built as **microservices and serverless functions** on
**Java 21 + Quarkus**, aligned with **The Twelve-Factor App**, documented under the ISO
architecture and quality standards, and aligned with the domain standards of the
vertical (see the traceability matrix). Everything is *aligned to* the cited standards —
**no certification, conformance or regulatory-compliance claim is made**. All data is
**synthetic**.

> Reference implementation: small, complete and runnable. The rules that matter in this
> domain are enforced in code and covered by tests — start with the invariant:
> credits issued against an ecosystem asset never exceed its verified service capacity — an issuance that would overrun it is rejected (422) and records nothing.

## Standards alignment

| Standard / framework | What it governs here | Where |
|---|---|---|
| ISO/IEC/IEEE **42010** | Architecture description | `docs/architecture/architecture-description.md` |
| ISO/IEC **25010** (SQuaRE) | Product quality model and scenarios | `docs/architecture/quality-attributes.md` |
| ISO/IEC/IEEE **15288 & 12207** | Life-cycle processes | `docs/lifecycle/software-lifecycle.md` |
| ISO/IEC **27001 & 27002** | Information security | `docs/security/` |
| **The Twelve-Factor App** | Cloud-native methodology | `docs/architecture/twelve-factor.md` |
| UN SEEA EA | Ecosystem accounting (extent, condition, services) | `docs/domain/domain-standards.md` mapping |
| ISO 14008 | Monetary valuation of environmental impacts | Valuation rules posture |
| ISO 14001 | Environmental management | Governance and life-cycle docs |
| TNFD | Nature-related financial disclosures | Boundary: disclosures consume this core's accounts |
| ISO 14064 | GHG quantification (carbon service) | Carbon service-type accounting |

Full traceability: `docs/compliance/standards-traceability.md`.

## The platform at a glance

### Microservices (Quarkus, `services/`)

| Service | Port | Alignment | Responsibility |
|---|---|---|---|
| `ecosystem-asset-service` | 8081 | UN SEEA EA ecosystem extent accounts | Registry of ecosystem assets under stewardship (SEEA EA extent accounts). |
| `service-credit-service` | 8082 | UN SEEA EA ecosystem services accounts | Verified service-capacity accounts per asset and the credit issuances against them. |
| `valuation-rules-service` | 8083 | ISO 14008 monetary valuation posture | Decision rules binding service types to value and uncertainty envelopes. |
| `monitoring-station-service` | 8084 | UN SEEA EA condition accounts (observation sources) | Registry of monitoring stations feeding ecosystem condition data. |
| `notification-service` | 8085 | Platform support (shared pattern) | Channel-agnostic dispatch of domain notifications from events. |

### Serverless functions (Quarkus Funqy + Knative, `functions/`)

| Function | Trigger | Responsibility |
|---|---|---|
| `ecosystem-health-scorer` | Event / HTTP | Scores degradation signals (deforestation, pollution, invasive pressure) with an auditable rationale. |
| `service-delivery-evaluator` | Event / HTTP | Compares degradation and service-index measurements against asset targets; one finding per breach. |

## Design principles

1. **One bounded context per service**, identical hexagonal layout everywhere.
2. **The domain invariant is code, not prose** — credits issued against an ecosystem asset never exceed its verified service capacity — an issuance that would overrun it is rejected (422) and records nothing.
3. **Privacy by design** — Land and asset data can be location-sensitive: precise coordinates are never required by any contract; assets are identified by internal UUIDs and region codes only.
4. **Database per service, API-first (OpenAPI), events as CloudEvents.**
5. **Adapt by configuration** — rules, thresholds, tariffs/rates and providers are data
   and configuration, not code forks (see *Adaptation guide*).

## Quick start

Prerequisites: Java 21, Maven 3.9+, optionally Docker/Podman and Knative/OpenShift.

```bash
mvn verify                                    # build and test everything
mvn -pl services/ecosystem-asset-service quarkus:dev
docker compose up -d                          # local PostgreSQL, Kafka, Keycloak
```

## Deployment

- **Kubernetes / OpenShift**: `kubectl apply -k deploy/kubernetes/`.
- **Serverless**: Knative Services in `deploy/knative/` (scale to zero).
- **Local**: `docker-compose.yml` for dependencies; services run with `quarkus:dev`.

## Adaptation guide

| You want to… | Change only… |
|---|---|
| Adjust domain rules and thresholds | Rule data in the `valuation-rules-service` context and function configuration |
| Change the guarded capacity semantics | The `service-credit-service` context (invariant stays: never negative) |
| Plug your IdP / providers | OIDC at the gateway; provider adapters at the boundary |
| Replace in-memory persistence | The repository port's JPA adapter on the service's own database (`DB_KIND`, `DB_URL`) |

## Documentation map

- Architecture (42010): `docs/architecture/` · Decisions: `docs/adr/`
- Domain standards: `docs/domain/domain-standards.md`
- Quality (25010), Security (27001/27002), Life cycle (15288/12207), API, Data, Operations: under `docs/`

## License

Apache-2.0 — see `LICENSE`.
