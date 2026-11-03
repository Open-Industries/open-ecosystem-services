package org.openindustries.ecosystem.asset;

import jakarta.validation.constraints.NotBlank;
import jakarta.validation.constraints.NotNull;
import jakarta.validation.constraints.Pattern;
import jakarta.validation.constraints.Positive;

public record CreateEcosystemAssetRequest(
        @NotBlank String name,
        @NotBlank @Pattern(regexp = "WETLAND|FOREST|MANGROVE|GRASSLAND|RIVER") String assetType,
        @Positive int hectares,
        @NotBlank String regionCode) {
}
