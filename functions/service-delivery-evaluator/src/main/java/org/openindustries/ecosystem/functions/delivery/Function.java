package org.openindustries.ecosystem.functions.delivery;

import io.quarkus.funqy.Funq;

import java.util.ArrayList;
import java.util.List;

/**
 * Serverless function (Knative / OpenShift Serverless).
 * Compares degradation and service-index measurements against asset targets; one finding per breach.
 * Compares measured values with targets and returns a verdict with one finding per
 * breach — machine-checkable assurance evidence.
 */
public class Function {

    public record Kpi(String assetId, int measuredDegradationPct, int targetDegradationPct,
                      double measuredServiceIndexPct, double targetServiceIndexPct) {
    }

    public record Verdict(String assetId, boolean compliant, List<String> findings) {
    }

    @Funq
    public Verdict serviceDeliveryEvaluator(Kpi kpi) {
        List<String> findings = new ArrayList<>();
        if (kpi.measuredDegradationPct() > kpi.targetDegradationPct()) {
            findings.add("degradationPct " + kpi.measuredDegradationPct() + " exceeds target " + kpi.targetDegradationPct());
        }
        if (kpi.measuredServiceIndexPct() < kpi.targetServiceIndexPct()) {
            findings.add("serviceIndexPct " + kpi.measuredServiceIndexPct() + " below target " + kpi.targetServiceIndexPct());
        }
        return new Verdict(kpi.assetId(), findings.isEmpty(), findings);
    }
}
