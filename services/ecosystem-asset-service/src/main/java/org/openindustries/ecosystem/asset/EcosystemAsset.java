package org.openindustries.ecosystem.asset;

import java.time.Instant;
import java.util.UUID;

/**
 * EcosystemAsset aggregate — Registry of ecosystem assets under stewardship (SEEA EA extent accounts).
 *
 * Privacy by design (ISO/IEC 27002 A.5.34): secret identifiers are validated at
 * creation and NEVER persisted or returned in clear — only a masked form is kept.
 * All identifiers in this reference implementation are synthetic.
 */
public class EcosystemAsset {
    public UUID id;
        public String name;
        public String assetType;
        public int hectares;
        public String regionCode;
    public String status;
    public Instant createdAt;

    public static EcosystemAsset create(String name, String assetType, int hectares, String regionCode) {
        EcosystemAsset e = new EcosystemAsset();
        e.id = UUID.randomUUID();
        e.name = name;
        e.assetType = assetType;
        e.hectares = hectares;
        e.regionCode = regionCode;
        e.status = "ACTIVE";
        e.createdAt = Instant.now();
        return e;
    }
}
