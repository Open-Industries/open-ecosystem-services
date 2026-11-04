package org.openindustries.ecosystem.credits;

import jakarta.enterprise.context.ApplicationScoped;
import jakarta.inject.Inject;
import java.util.List;
import java.util.UUID;

@ApplicationScoped
public class CreditsService {

    @Inject
    CreditsRepository repository;

    public CreditAccount open(OpenCreditAccountRequest request) {
        return repository.saveAccount(CreditAccount.open(request.assetId(), request.capacityUnits()));
    }

    public CreditAccount get(UUID id) {
        return repository.findAccount(id)
                .orElseThrow(() -> new DomainException(404, "CreditAccount not found: " + id));
    }

    /**
     * Draws from the account. Invariant: the balance/capacity never goes negative —
     * an over-drawing operation is rejected (422) and nothing is recorded.
     */
    public CreditIssuance consume(UUID accountId, CreateCreditIssuanceRequest request) {
        CreditAccount account = get(accountId);
        if (!"ACTIVE".equals(account.status)) {
            throw new DomainException(409, "CreditAccount is not ACTIVE");
        }
        if (request.units() > account.capacityUnits) {
            throw new DomainException(422, "Insufficient capacity: requested " + request.units()
                    + " exceeds available " + account.capacityUnits);
        }
        account.capacityUnits -= request.units();
        repository.saveAccount(account);
        return repository.saveRecord(CreditIssuance.of(accountId, request.serviceType(), request.units(), account.capacityUnits));
    }

    public List<CreditIssuance> recordsOf(UUID accountId) {
        get(accountId);
        return repository.recordsOf(accountId);
    }
}
