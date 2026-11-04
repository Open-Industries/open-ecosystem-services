package org.openindustries.ecosystem.credits;

import jakarta.validation.constraints.NotBlank;
import jakarta.validation.constraints.Pattern;
import jakarta.validation.constraints.Positive;

public record CreateCreditIssuanceRequest(
        @NotBlank @Pattern(regexp = "CARBON|WATER|BIODIVERSITY|POLLINATION") String serviceType,
        @Positive long units) {
}
