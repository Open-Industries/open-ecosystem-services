# Monitoring Station Service

Bounded context: **Monitoring Stations** — Registry of monitoring stations feeding ecosystem condition data.

Alignment: UN SEEA EA condition accounts (observation sources). Lifecycle transitions take no body. All data is synthetic.

| Item | Value |
|---|---|
| Port | `8084` (`HTTP_PORT`) |
| Base path | `/api/v1/monitoring-stations` |
| Owns data | MonitoringStation records |

Run locally: `mvn -pl services/monitoring-station-service quarkus:dev`.
