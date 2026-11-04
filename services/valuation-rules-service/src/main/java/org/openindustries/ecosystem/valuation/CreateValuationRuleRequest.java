package org.openindustries.ecosystem.valuation;

import jakarta.validation.constraints.Max;
import jakarta.validation.constraints.Min;
import jakarta.validation.constraints.NotBlank;
import jakarta.validation.constraints.Pattern;
import jakarta.validation.constraints.Positive;

public record CreateValuationRuleRequest(
        @NotBlank String name,
        @NotBlank @Pattern(regexp = "CARBON|WATER|BIODIVERSITY|POLLINATION") String serviceType,
        @Positive long maxValuePerHectareMinor,
        @Positive int maxUncertaintyPct,
        @Min(1) @Max(100) int priority) {
}
