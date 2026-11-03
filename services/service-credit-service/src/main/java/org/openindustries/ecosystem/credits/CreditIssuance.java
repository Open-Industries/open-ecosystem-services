package org.openindustries.ecosystem.credits;

import java.time.Instant;
import java.util.UUID;

/** One guarded consumption record — the auditable evidence of a draw on a CreditAccount. */
public class CreditIssuance {
    public UUID id;
    public UUID accountId;
    public String serviceType;
    public long units;
    public long remainingAfter;
    public Instant recordedAt;

    public static CreditIssuance of(UUID accountId, String serviceType, long units, long remainingAfter) {
        CreditIssuance r = new CreditIssuance();
        r.id = UUID.randomUUID();
        r.accountId = accountId;
        r.serviceType = serviceType;
        r.units = units;
        r.remainingAfter = remainingAfter;
        r.recordedAt = Instant.now();
        return r;
    }
}
