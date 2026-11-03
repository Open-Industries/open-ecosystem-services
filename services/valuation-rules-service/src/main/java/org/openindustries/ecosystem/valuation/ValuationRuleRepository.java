package org.openindustries.ecosystem.valuation;

import jakarta.enterprise.context.ApplicationScoped;
import java.util.Comparator;
import java.util.List;
import java.util.Optional;
import java.util.UUID;
import java.util.concurrent.ConcurrentHashMap;

/** Outbound port; default in-memory adapter, production adapter is JPA on this service's own database. */
@ApplicationScoped
public class ValuationRuleRepository {
    private final ConcurrentHashMap<UUID, ValuationRule> store = new ConcurrentHashMap<>();

    public ValuationRule save(ValuationRule rule) {
        store.put(rule.id, rule);
        return rule;
    }

    public Optional<ValuationRule> findById(UUID id) {
        return Optional.ofNullable(store.get(id));
    }

    public List<ValuationRule> findAll() {
        return List.copyOf(store.values());
    }

    /** Active rules for a scope, highest priority (lowest number) first. */
    public List<ValuationRule> activeFor(String serviceType) {
        return store.values().stream()
                .filter(r -> r.serviceType.equals(serviceType) && "ACTIVE".equals(r.status))
                .sorted(Comparator.comparingInt(r -> r.priority))
                .toList();
    }
}
