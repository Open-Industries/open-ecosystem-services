# Service Delivery Evaluator (serverless function)

Quarkus Funqy function on Knative (scale to zero). Compares degradation and service-index measurements against asset targets; one finding per breach.

Invoke: `POST /serviceDeliveryEvaluator` with the KPI payload.
