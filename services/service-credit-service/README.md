# Service Credit Service

Bounded context: **Service Credits** — Verified service-capacity accounts per asset and the credit issuances against them.

Alignment: UN SEEA EA ecosystem services accounts. Core invariant, enforced in code and tests: **the guarded
balance/capacity never goes negative** — over-drawing operations are rejected with `422`
and nothing is recorded. All data is synthetic.

| Item | Value |
|---|---|
| Port | `8082` (`HTTP_PORT`) |
| Base path | `/api/v1/credit-accounts` |
| Owns data | CreditAccount accounts, CreditIssuance records |
| Publishes | `credits.consumed`, `credits.insufficient` |

Run locally: `mvn -pl services/service-credit-service quarkus:dev`.
