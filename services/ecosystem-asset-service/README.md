# Ecosystem Asset Service

Bounded context: **Ecosystem Assets** — Registry of ecosystem assets under stewardship (SEEA EA extent accounts).

Alignment: UN SEEA EA ecosystem extent accounts. Lifecycle transitions take no body; secret identifiers,
where the domain has them, are masked at creation and never stored or returned in clear.
All data is synthetic.

| Item | Value |
|---|---|
| Port | `8081` (`HTTP_PORT`) |
| Base path | `/api/v1/ecosystem-assets` |
| Owns data | EcosystemAsset records |
| Publishes | `asset.created`, lifecycle events (see `contracts/asyncapi.yaml`) |

Run locally: `mvn -pl services/ecosystem-asset-service quarkus:dev`.
