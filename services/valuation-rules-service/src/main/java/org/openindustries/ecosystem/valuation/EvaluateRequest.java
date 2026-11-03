package org.openindustries.ecosystem.valuation;

import jakarta.validation.constraints.NotBlank;
import jakarta.validation.constraints.Pattern;
import jakarta.validation.constraints.Positive;

public record EvaluateRequest(
        @NotBlank @Pattern(regexp = "CARBON|WATER|BIODIVERSITY|POLLINATION") String serviceType,
        @Positive int requiredMaxUncertaintyPct,
        @Positive long requiredMaxValuePerHectareMinor) {
}
