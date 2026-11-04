package org.openindustries.ecosystem.credits;

import jakarta.validation.constraints.NotBlank;
import jakarta.validation.constraints.NotNull;
import jakarta.validation.constraints.Pattern;
import jakarta.validation.constraints.Positive;
import jakarta.validation.constraints.PositiveOrZero;

public record OpenCreditAccountRequest(
        @NotNull java.util.UUID assetId,
        @PositiveOrZero long capacityUnits) {
}
