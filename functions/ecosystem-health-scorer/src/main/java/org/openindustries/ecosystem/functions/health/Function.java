package org.openindustries.ecosystem.functions.health;

import io.quarkus.funqy.Funq;

/**
 * Serverless function (Knative / OpenShift Serverless).
 * Scores degradation signals (deforestation, pollution, invasive pressure) with an auditable rationale.
 * Transparent, auditable rules: NORMAL (< 35), WARNING (35-69), CRITICAL (>= 70);
 * the rationale is returned with the score for audit (ISO/IEC 27002 A.8.15).
 * Stateless by design (12-Factor VI): all state arrives in the request.
 */
public class Function {

    public record Event(String assetId, int deforestationIndex, int pollutionIndex, int invasiveIndex) {
    }

    public record Score(String assetId, int score, String decision, String rationale) {
    }

    @Funq
    public Score ecosystemHealthScorer(Event event) {
        int score = 0;
        StringBuilder why = new StringBuilder();
        if (event.deforestationIndex() >= 70) {
            score += 40;
            why.append("deforestationIndex >= 70; ");
        }
        if (event.pollutionIndex() >= 60) {
            score += 30;
            why.append("pollutionIndex >= 60; ");
        }
        if (event.invasiveIndex() >= 50) {
            score += 30;
            why.append("invasiveIndex >= 50; ");
        }
        String decision = score >= 70 ? "CRITICAL" : score >= 35 ? "WARNING" : "NORMAL";
        return new Score(event.assetId(), score, decision,
                why.length() == 0 ? "no signals" : why.toString().trim());
    }
}
