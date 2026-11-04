package org.openindustries.ecosystem.station;

import jakarta.validation.constraints.NotBlank;
import jakarta.validation.constraints.NotNull;
import jakarta.validation.constraints.Pattern;
import jakarta.validation.constraints.Positive;

public record CreateMonitoringStationRequest(
        @NotBlank String name,
        @NotBlank @Pattern(regexp = "FLUX_TOWER|WATER_GAUGE|CAMERA_TRAP|SOIL_PROBE") String stationType,
        @Positive int readingsPerDay) {
}
