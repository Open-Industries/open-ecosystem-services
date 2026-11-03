package org.openindustries.ecosystem.valuation;

import java.util.UUID;

/**
 * ValuationRule aggregate — a decision rule of this core: for one serviceType, it guarantees a
 * capacity (maxValuePerHectareMinor) and a threshold (maxUncertaintyPct), with a precedence priority.
 * Only ACTIVE rules decide; every decision names its rule and its reason (audit).
 */
public class ValuationRule {
    public UUID id;
    public String name;
    public String serviceType;
    public long maxValuePerHectareMinor;
    public int maxUncertaintyPct;
    public int priority;
    public String status;

    public static ValuationRule create(String name, String serviceType, long maxValuePerHectareMinor, int maxUncertaintyPct, int priority) {
        ValuationRule r = new ValuationRule();
        r.id = UUID.randomUUID();
        r.name = name;
        r.serviceType = serviceType;
        r.maxValuePerHectareMinor = maxValuePerHectareMinor;
        r.maxUncertaintyPct = maxUncertaintyPct;
        r.priority = priority;
        r.status = "DRAFT";
        return r;
    }
}
