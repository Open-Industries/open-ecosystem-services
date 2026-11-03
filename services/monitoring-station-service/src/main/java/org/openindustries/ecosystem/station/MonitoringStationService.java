package org.openindustries.ecosystem.station;

import jakarta.enterprise.context.ApplicationScoped;
import jakarta.inject.Inject;
import java.util.List;
import java.util.UUID;

@ApplicationScoped
public class MonitoringStationService {

    @Inject
    MonitoringStationRepository repository;

    public MonitoringStation create(CreateMonitoringStationRequest request) {
        return repository.save(MonitoringStation.create(request.name(), request.stationType(), request.readingsPerDay()));
    }

    public MonitoringStation get(UUID id) {
        return repository.findById(id)
                .orElseThrow(() -> new DomainException(404, "MonitoringStation not found: " + id));
    }

    public List<MonitoringStation> list() {
        return repository.findAll();
    }

    public MonitoringStation activate(UUID id) {
        MonitoringStation entity = get(id);
        if (!"INSTALLED".equals(entity.status)) {
            throw new DomainException(409, "Only a INSTALLED monitoringstation can transition via activate");
        }
        entity.status = "ACTIVE";
        return repository.save(entity);
    }
}
