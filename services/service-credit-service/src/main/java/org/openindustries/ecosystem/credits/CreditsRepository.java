package org.openindustries.ecosystem.credits;

import jakarta.enterprise.context.ApplicationScoped;
import java.util.List;
import java.util.Optional;
import java.util.UUID;
import java.util.concurrent.ConcurrentHashMap;
import java.util.concurrent.CopyOnWriteArrayList;

/** Outbound port; default in-memory adapter, production adapter is JPA on this service's own database. */
@ApplicationScoped
public class CreditsRepository {
    private final ConcurrentHashMap<UUID, CreditAccount> accounts = new ConcurrentHashMap<>();
    private final CopyOnWriteArrayList<CreditIssuance> records = new CopyOnWriteArrayList<>();

    public CreditAccount saveAccount(CreditAccount account) {
        accounts.put(account.id, account);
        return account;
    }

    public Optional<CreditAccount> findAccount(UUID id) {
        return Optional.ofNullable(accounts.get(id));
    }

    public CreditIssuance saveRecord(CreditIssuance record) {
        records.add(record);
        return record;
    }

    public List<CreditIssuance> recordsOf(UUID accountId) {
        return records.stream().filter(r -> r.accountId.equals(accountId)).toList();
    }
}
