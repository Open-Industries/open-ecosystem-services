package org.openindustries.ecosystem.station;

import jakarta.enterprise.context.ApplicationScoped;
import java.util.List;
import java.util.Optional;
import java.util.UUID;
import java.util.concurrent.ConcurrentHashMap;

/** Outbound port; default in-memory adapter, production adapter is JPA on this service's own database. */
@ApplicationScoped
public class MonitoringStationRepository {
    private final ConcurrentHashMap<UUID, MonitoringStation> store = new ConcurrentHashMap<>();

    public MonitoringStation save(MonitoringStation entity) {
        store.put(entity.id, entity);
        return entity;
    }

    public Optional<MonitoringStation> findById(UUID id) {
        return Optional.ofNullable(store.get(id));
    }

    public List<MonitoringStation> findAll() {
        return List.copyOf(store.values());
    }
}
