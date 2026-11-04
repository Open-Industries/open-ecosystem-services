package org.openindustries.ecosystem.credits;

import java.util.UUID;

/**
 * CreditAccount aggregate — the guarded resource account of this core: a balance/capacity that
 * consumption operations draw from. Core invariant (ADR-0004): it NEVER goes negative;
 * an operation that would overdraw it is rejected and records nothing.
 */
public class CreditAccount {
    public UUID id;
        public java.util.UUID assetId;
    public long capacityUnits;
    public String status;

    public static CreditAccount open(java.util.UUID assetId, long capacityUnits) {
        CreditAccount a = new CreditAccount();
        a.id = UUID.randomUUID();
        a.assetId = assetId;
        a.capacityUnits = capacityUnits;
        a.status = "ACTIVE";
        return a;
    }
}
