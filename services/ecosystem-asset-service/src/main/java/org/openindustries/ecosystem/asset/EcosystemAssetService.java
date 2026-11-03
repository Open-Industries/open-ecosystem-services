package org.openindustries.ecosystem.asset;

import jakarta.enterprise.context.ApplicationScoped;
import jakarta.inject.Inject;
import java.util.List;
import java.util.UUID;

@ApplicationScoped
public class EcosystemAssetService {

    @Inject
    EcosystemAssetRepository repository;

    public EcosystemAsset create(CreateEcosystemAssetRequest request) {
        return repository.save(EcosystemAsset.create(request.name(), request.assetType(), request.hectares(), request.regionCode()));
    }

    public EcosystemAsset get(UUID id) {
        return repository.findById(id)
                .orElseThrow(() -> new DomainException(404, "EcosystemAsset not found: " + id));
    }

    public List<EcosystemAsset> list() {
        return repository.findAll();
    }

    public EcosystemAsset protect(UUID id) {
        EcosystemAsset entity = get(id);
        if ("PROTECTED".equals(entity.status)) {
            throw new DomainException(409, "EcosystemAsset already PROTECTED");
        }
        entity.status = "PROTECTED";
        return repository.save(entity);
    }

    public EcosystemAsset restore(UUID id) {
        EcosystemAsset entity = get(id);
        if ("UNDER_RESTORATION".equals(entity.status)) {
            throw new DomainException(409, "EcosystemAsset already UNDER_RESTORATION");
        }
        entity.status = "UNDER_RESTORATION";
        return repository.save(entity);
    }
}
